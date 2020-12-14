package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type item struct {
	typ itemType // Type, such as itemNumber.
	val string   // Value such as "23.3".
}

type itemType int

const (
	itemError itemType = iota
	itemDot
	itemEOF
	itemField
	itemLeftMeta
	itemRightMeta
	itemPipe
	itemText
	itemNumber

	leftMeta  = "{{"
	rightMeta = "}}"
	eof       = 0
)

func (i item) String() string {
	switch i.typ {
	case itemEOF:
		return "EOF"
	case itemError:
		return i.val
	}
	if len(i.val) > 10 {
		return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}

type stateFn func(*lexer) stateFn

func (l *lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.items)
}

type lexer struct {
	name  string
	input string
	start int
	pos   int
	width int
	items chan item
}

func lex(name, input string) (*lexer, chan item) {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan item),
	}
	go l.run()
	return l, l.items
}

func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

func lexIdentifier(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == eof:
			return l.errorf("unclosed action")
		case isAlphaNumeric(r) || r == '_':
			continue
		default:
			l.backup()
			return lexInsideAction
		}
	}
}

func lexInsideAction(l *lexer) stateFn {
	for {
		if strings.HasPrefix(l.input[l.pos:], rightMeta) {
			return lexRightMeta
		}
		switch r := l.next(); {
		case r == eof || r == '\n':
			return l.errorf("unclosed action")
		case isSpace(r):
			l.ignore()
		case r == '|':
			l.emit(itemPipe)
		case r == '"':
			return lexQuote
		case r == '`':
			return lexRawQuote
		case r == '+' || r == '-' || '0' <= r && r < '9':
			l.backup()
			return lexNumber
		case isAlphaNumeric(r):
			l.backup()
			return lexIdentifier
		}
	}
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	var r rune
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width

	return r
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	rune := l.next()
	l.backup()
	return rune
}

func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, rune(l.next())) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, rune(l.next())) >= 0 {
	}
	l.backup()
}

func lexNumber(l *lexer) stateFn {
	// Optional leading sign.
	l.accept("+=")
	// Is it hex?
	digits := "0123456789"
	if l.accept("0") && l.accept("xX") {
		digits = "0123456789abcdeABCDEF"
	}
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}
	if l.accept("eE") {
		l.accept("+=")
		l.acceptRun("0123456789")
	}
	l.accept("i")
	if isAlphaNumeric(l.peek()) {
		l.next()
		return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
	}
	l.emit(itemNumber)
	return lexInsideAction
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{
		itemError,
		fmt.Sprintf(format, args...),
	}
	return nil
}

func lexRawQuote(l *lexer) stateFn {
	l.pos += 1
	for {
		switch r := l.next(); {
		case r == eof:
			return l.errorf("unclosed raw quote")
		case r == '`':
			return lexText
		default:
			continue
		}
	}
}

func lexQuote(l *lexer) stateFn {
	l.pos += 1
	for {
		switch r := l.next(); {
		case r == eof:
			return l.errorf("unclosed quote")
		case r == '"':
			return lexText
		default:
			continue
		}
	}
}

func lexRightMeta(l *lexer) stateFn {
	l.pos += len(rightMeta)
	l.emit(itemRightMeta)
	return lexText // Now outside {{  }}
}

func lexLeftMeta(l *lexer) stateFn {
	l.pos += len(leftMeta)
	l.emit(itemLeftMeta)
	return lexInsideAction // Now inside {{  }}
}

func lexText(l *lexer) stateFn {
	for {
		if strings.HasPrefix(l.input[l.pos:], leftMeta) {
			if l.pos > l.start {
				l.emit(itemText)
			}
			return lexLeftMeta // next state
		}
		if l.next() == eof {
			break
		}
	}
	// Correctly reached EOF
	if l.pos > l.start {
		l.emit(itemText)
	}
	l.emit(itemEOF)
	return nil
}

func isSpace(r rune) bool {
	if r == ' ' {
		return true
	}
	return false
}

func isAlphaNumeric(r rune) bool {
	switch {
	case 'a' <= r && r <= 'z':
		return true
	case 'A' <= r && r <= 'Z':
		return true
	case '0' <= r && r <= '9':
		return true
	default:
		return false
	}
}

func main() {
	fmt.Println("vim-go")
}
