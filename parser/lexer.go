package parser

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type item struct {
	typ  itemType
	pos  int
	val  string
	line int
}

func (i item) String() string {
	switch {
	case i.typ == itemEOF:
		return "EOF"
	case i.typ == itemError:
		return i.val
	case i.typ > itemKeyword:
		return fmt.Sprintf("<%s>", i.val)
	case len(i.val) > 10:
		return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}

// itemType identifies the type of lex items.
type itemType int

const (
	itemError itemType = iota
	itemEOF
	itemSpace
	itemLeftDelim
	itemRightDelim
	itemColon
	itemComma
	itemString
	itemKeyword
	itemFile
	itemDeps
	itemCompiler
	itemBinary
	itemName
	itemStart
	itemLdFlags
	itemIncludes
	itemOtherFlags
)

var key = map[string]itemType{
	"file":     itemFile,
	"deps":     itemDeps,
	"$":        itemCompiler,
	"binary":   itemBinary,
	"name":     itemName,
	"start":    itemStart,
	"ldflags":  itemLdFlags,
	"includes": itemIncludes,
	"others":   itemOtherFlags,
}

const eof = -1

type lexer struct {
	name       string
	input      string
	leftDelim  string
	rightDelim string
	pos        int
	start      int
	width      int
	items      chan item
	parenDepth int
	line       int
}

type stateFn func(*lexer) stateFn

func (self *lexer) next() rune {
	if self.pos >= len(self.input) {
		self.width = 0
		return eof
	}

	nextRune, runeWidth := utf8.DecodeRuneInString(self.input[self.pos:])
	self.width = runeWidth
	self.pos += self.width
	if nextRune == '\n' {
		self.line++
	}
	return nextRune
}

func (self *lexer) peek() rune {
	nextRune := self.next()
	self.backup()
	return nextRune
}

func (self *lexer) backup() {
	self.pos -= self.width
	if self.width == 1 && self.input[self.pos] == '\n' {
		self.line--
	}
}

func (self *lexer) emit(t itemType) {
	self.items <- item{t, self.start, self.input[self.start:self.pos], self.line}

	switch t {
	case itemString, itemLeftDelim, itemRightDelim:
		self.line += strings.Count(self.input[self.start:self.pos], "\n")
	}
	self.start = self.pos
}

func (self *lexer) ignore() {
	self.line += strings.Count(self.input[self.start:self.pos], "\n")
	self.start = self.pos
}
