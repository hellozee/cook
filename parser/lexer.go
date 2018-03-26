package parser

import "fmt"

type item struct {
	typ  itemType
	pos  Pos
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
	itemLeftBrackets
	itemRightBrackets
	itemColon
	itemComma
	itemKeyword
	itemFile
	itemDeps
	itemString
	itemCompiler
	itemBinary
	itemName
	itemStart
	itemLdFlags
	itemIncludes
	itemOtherFlags
)

// Pos represents a byte position in the original input text from which
// this template was parsed.
type Pos int

func (p Pos) Position() Pos {
	return p
}
