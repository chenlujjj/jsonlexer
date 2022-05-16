package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type item struct {
	typ itemType
	val string
}

func (i item) String() string {
	switch {
	case i.typ == itemEOF:
		return "EOF"
	case i.typ == itemError:
		return i.val
	default:
		return fmt.Sprintf("%q", i.val)
	}
}

type itemType int

const (
	itemError itemType = iota
	itemDoubleQuote
	itemLeftBrace    // {
	itemRightBrace   // }
	itemLeftBracket  // [
	itemRightBracket // ]
	itemColon        // :
	itemComma        // ,
	itemSpace
	itemNumber
	itemString
	itemObject
	itemArray
	itemTrue  // true
	itemFalse // false
	itemNull  // null
	itemValue
	itemEOF
)

type stateFn func(*lexer) stateFn

type lexer struct {
	name  string // the name of the input; used only for error reports
	input string // the string being scanned
	pos   int    // current position in the input
	start int    // start position of this item
	width int    // width of last rune read from input 也就是rune的字节数
	items []item // slice or channel ?

	inObjectDepth int // 记录在object中的深度，据此来判断是否应该结束json
	inArrayDepth  int
}

func (l *lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
}

func (l *lexer) emit(t itemType) {
	l.items = append(l.items, item{t, l.input[l.start:l.pos]})
	l.start = l.pos
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}

func (l *lexer) forward(w int) {
	l.width = w
	l.pos += w
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) errorf(format string, args ...any) stateFn {
	l.items = append(l.items, item{itemError, fmt.Sprintf(format, args...)})
	return nil
}

const (
	eof        = -1
	spaceChars = " \t\r\n" // These are the space characters defined by Go itself.

	leftBrace    = "{"
	rightBrace   = "}"
	leftBracket  = "["
	rightBracket = "]"
	comma        = ","
	colon        = ':'
	doubleQuote  = `"`
)

// leftTrimLength returns the length of the spaces at the beginning of the string.
func leftTrimLength(s string) int {
	return len(s) - len(strings.TrimLeft(s, spaceChars))
}

// 移除开头的space
func lexText(l *lexer) stateFn {
	l.pos += leftTrimLength(l.input[l.pos:])
	l.start = l.pos

	return lexValue
}

func lexValue(l *lexer) stateFn {
	r := l.next()
	switch {
	case isSpace(r):
		l.ignore()
		return lexValue
	case r == '{':
		l.emit(itemLeftBrace)
		l.inObjectDepth += 1
		return lexInside
	case r == '[':
		l.emit(itemLeftBracket)
		l.inArrayDepth += 1
		return lexInside
	case r == eof:
		return nil
	case r == '"':
		return lexQuote
	default:
		l.backup()
		if strings.HasPrefix(l.input[l.pos:], "true") {
			return lexTrue
		} else if strings.HasPrefix(l.input[l.pos:], "false") {
			return lexFalse
		} else if strings.HasPrefix(l.input[l.pos:], "null") {
			return lexNull
		} else {
			// TODO 单个数字
			return l.errorf("unexpected character %#U", r)
		}
	}
}

// lexQuote scans a quoted string.
// TODO 先不考虑字符串里有引号的情况
func lexQuote(l *lexer) stateFn {
Loop:
	for {
		switch l.next() {
		case eof:
			return l.errorf("unterminated quoted string")
		case '"':
			break Loop
		}
	}
	l.emit(itemString)
	if l.inObjectDepth > 0 || l.inArrayDepth > 0 {
		return lexInside
	}
	return lexValue
}

func lexInside(l *lexer) stateFn {
	r := l.next()
	switch {
	case isSpace(r):
		l.ignore()
		return lexInside
	case r == '"':
		return lexQuote
	case r == ':':
		l.emit(itemColon)
		return lexInside
	case r == ',':
		l.emit(itemComma)
		return lexInside
	case r == eof:
		return l.errorf("unexpected eof")
	case r == '}':
		l.emit(itemRightBrace)
		l.inObjectDepth -= 1
		if l.inObjectDepth == 0 && l.inArrayDepth == 0 {
			// 结束
			return nil
		} else if l.inObjectDepth < 0 {
			return l.errorf("unexpected right brace")
		}
	case r == '{':
		l.emit(itemLeftBrace)
		l.inObjectDepth += 1
		return lexInside
	case r == '[':
		l.emit(itemLeftBracket)
		l.inArrayDepth += 1
		return lexInside
	case r == ']':
		l.emit(itemRightBracket)
		l.inArrayDepth -= 1
		if l.inArrayDepth == 0 && l.inObjectDepth == 0 {
			// 结束
			return nil
		} else if l.inArrayDepth < 0 {
			return l.errorf("unexpected right bracket")
		}
	default:
		l.backup()
		if strings.HasPrefix(l.input[l.pos:], "true") {
			return lexTrue
		} else if strings.HasPrefix(l.input[l.pos:], "false") {
			return lexFalse
		} else if strings.HasPrefix(l.input[l.pos:], "null") {
			return lexNull
		} else {
			// TODO 单个数字
			return l.errorf("unexpected character %#U", r)
		}
	}
	return lexInside
}

func lexTrue(l *lexer) stateFn {
	l.forward(len("true"))
	l.emit(itemTrue)
	if l.inObjectDepth > 0 || l.inArrayDepth > 0 {
		return lexInside
	}
	return lexValue
}

func lexFalse(l *lexer) stateFn {
	l.forward(len("false"))
	l.emit(itemFalse)
	if l.inObjectDepth > 0 || l.inArrayDepth > 0 {
		return lexInside
	}
	return lexValue
}

func lexNull(l *lexer) stateFn {
	l.forward(len("null"))
	l.emit(itemNull)
	if l.inObjectDepth > 0 || l.inArrayDepth > 0 {
		return lexInside
	}
	return lexValue
}


// 整数，小数，科学计数法
func lexNumber(l *lexer) stateFn {
	return nil
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r' || r == '\n'
}
