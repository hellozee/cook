package parser

import (
	"strings"
	"unicode/utf8"
)

//item  Data Structure for holding a single item
type item struct {
	typ  itemType
	pos  int
	val  string
	line int
}

//itemType  Substituting int with itemType for better readability
type itemType int

//itemEnum  Enum for holding the various item types
const (
	// Literal terminator symbols
	itemEOF itemType = iota
	itemNULL
	itemEquals
	itemSemicolon
	itemLeftBrace
	itemRightBrace
	itemString
	itemEntity
	itemlineBreak
	itemKeyWord
	itemBinary
	itemName
	itemStart
	itemLdFlags
	itemIncludes
	itemOthers
	itemFile
	itemDeps
)

//eof  Constant for determining the end of file
const eof = -1

//key  Map for holding various keywords with respect to item enum
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

//lexer Data Structure for holding the Lexer
type lexer struct {
	input string
	pos   int
	start int
	width int
	items []item
	line  int
}

//analyze  Function to perform lexical analysis on the stream
func (lex *lexer) analyze() {
	for lex.peek() != eof {
		lex.next()
		lex.isKeyword()
		if lex.peek() == eof {
			break
		}
		lex.isDelimiter()
	}

	lex.items = append(lex.items, item{typ: itemEOF, pos: lex.pos + 1,
		val: "EOF", line: lex.line})
}

//next  Function to shift the next rune in the stream
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

//peek  Function to check the next rune while mainting the current position
func (lex *lexer) peek() rune {
	nextRune := lex.next()
	lex.backup()
	return nextRune
}

//backup  Function to back one rune from the current position of the Lexer
func (lex *lexer) backup() {
	lex.pos -= lex.width

	if lex.width == 1 && lex.input[lex.pos] == '\n' {
		lex.line--
	}
}

/*isKeyword  Function to to determine whether the given character or
  the stream of characters is a keyword or not, if yes then
  the appropriate item is pushed into the lexer data structure */
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

/*isDelimiter  Function to to determine whether the given character or
  the stream of characters is a delimiter or not, if yes then
  the appropriate item is pushed into the lexer data structure */
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

	case '\n':
		if lex.items[len(lex.items)-1].typ == itemEquals {
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
		}
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
		if lex.items[len(lex.items)-1].typ != itemEquals {
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
}

//newLexer Function to create a new Lexer
func newLexer(file string) *lexer {
	lex := lexer{
		input: file,
		pos:   0,
		start: 0,
		width: len(file),
		line:  0,
	}

	return &lex
}

//isWhiteSpace Function to determine if the give character is a white space or not
func isWhiteSpace(currentRune rune) bool {
	return currentRune == ' ' || currentRune == '\t' || currentRune == '\r'
}
