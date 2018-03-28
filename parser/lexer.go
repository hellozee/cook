package main

import (
	"unicode/utf8"
)

type itemType int

const (
	// Literal terminator symbols
	itemEquals itemType = iota
	itemSemicolon
	itemLeftBrace
	itemrRightBrace
	//KeyWords
	itemEntity
	itemBinary
	itemName
	itemStart
	itemLdFlags
	itemIncludes
	itemOthers
	itemFile
	itemDeps
)

const eof = -1

var key = map[string]itemType{
	"entity":   itemEntity,
	"binary":   itemBinary,
	"name":     itemName,
	"start":    itemStart,
	"ldflags":  itemLdFlags,
	"includes": itemIncludes,
	"others":   itemOthers,
	"file":     itemFile,
	"deps":     itemDeps,
}

type lexer struct {
	input string
	pos   int
	start int
	width int
	line  int
}

func (lex *lexer) next() rune {
	if lex.pos >= len(lex.input) {
		return eof
	}

	nextRune, runeWidth := utf8.DecodeRuneInString(lex.input[lex.pos:])
	lex.width = runeWidth
	lex.pos += lex.width
	if isWhiteSpace(nextRune) {
		return lex.next()
	}
	if nextRune == '\n' {
		lex.line++
	}
	return nextRune
}

func (lex *lexer) peek() rune {
	nextRune := lex.next()
	lex.backup()
	return nextRune
}

func (lex *lexer) backup() {
	lex.pos -= lex.width

	if lex.width == 1 && lex.input[lex.pos] == '\n' {
		lex.line--
	}
}

func isWhiteSpace(currentRune rune) bool {
	return currentRune == ' ' || currentRune == '\t' || currentRune == '\r'
}
