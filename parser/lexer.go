package main

import (
	"strings"
	"unicode/utf8"
)

type item struct {
	typ  itemType
	pos  int
	val  string
	line int
}

type itemType int

const (
	// Literal terminator symbols
	itemEquals itemType = iota
	itemSemicolon
	itemLeftBrace
	itemRightBrace
	itemString
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
	items []item
	line  int
}

func (lex *lexer) analyze() {
	for lex.peek() != eof {
		lex.next()
		lex.isKeyword()
		if lex.peek() == eof {
			break
		}
		lex.isDelimiter()
	}
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

func (lex *lexer) isKeyword() {
	value := strings.TrimSpace(lex.input[lex.start:lex.pos])

	if typeOf, ok := key[value]; ok {
		tempItem := item{
			typ:  typeOf,
			pos:  lex.start,
			val:  value,
			line: lex.line,
		}

		lex.items = append(lex.items, tempItem)
		lex.start = lex.pos
	}
}

func (lex *lexer) isDelimiter() {
	switch lex.input[lex.pos] {
	case '{':
		tempItem := item{
			typ:  itemString,
			pos:  lex.start,
			val:  strings.TrimSpace(lex.input[lex.start:lex.pos]),
			line: lex.line,
		}
		lex.items = append(lex.items, tempItem)

		tempItem = item{
			typ:  itemLeftBrace,
			pos:  lex.pos,
			val:  "{",
			line: lex.line,
		}
		lex.items = append(lex.items, tempItem)

		lex.start = lex.pos + 1

	case '}':
		tempItem := item{
			typ:  itemRightBrace,
			pos:  lex.pos,
			val:  "}",
			line: lex.line,
		}
		lex.items = append(lex.items, tempItem)

		lex.start = lex.pos + 1

	case ';':
		tempItem := item{
			typ:  itemString,
			pos:  lex.start,
			val:  strings.TrimSpace(lex.input[lex.start:lex.pos]),
			line: lex.line,
		}
		lex.items = append(lex.items, tempItem)

		tempItem = item{
			typ:  itemSemicolon,
			pos:  lex.pos,
			val:  ";",
			line: lex.line,
		}
		lex.items = append(lex.items, tempItem)
		lex.start = lex.pos + 1

	case '=':
		tempItem := item{
			typ:  itemEquals,
			pos:  lex.pos,
			val:  "=",
			line: lex.line,
		}
		lex.items = append(lex.items, tempItem)
		lex.start = lex.pos + 1
	}
}

func isWhiteSpace(currentRune rune) bool {
	return currentRune == ' ' || currentRune == '\t' || currentRune == '\r'
}
