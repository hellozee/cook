package parser

import (
	"errors"
	"strings"
)

type compiler struct {
	binary     string
	name       string
	start      string
	ldFlags    string
	includes   string
	otherFlags string
}

type params struct {
	file string
	deps []string
}

type parser struct {
	input       []item
	pos         int
	currentItem item
	prevItem    item
	nextItem    item
}

func (par *parser) next() item {
	par.prevItem = par.currentItem
	par.currentItem = par.nextItem
	par.pos++
	par.nextItem = par.input[par.pos]
	return par.currentItem
}

func (par *parser) parse() error {

	isCompiler := false
	enityName := ""
	identifier := itemNULL
	params := ""

	for par.nextItem.typ != itemEOF {
		par.next()
		if par.currentItem.typ == itemEquals {
			if par.nextItem.typ != itemString {
				if par.prevItem.typ == itemEntity {
					return par.reportError("entity")
				}
				return par.reportError("parameter(s)")
			}
		}

		if par.currentItem.typ == itemSemicolon {
			if par.nextItem.typ < itemKeyWord &&
				par.nextItem.typ != itemRightBrace {
				return par.reportError("}")
			}

			if isCompiler {
				fillCompilerDetails(identifier, params)

			} else {
				fillFileDetails(enityName, identifier, params)
			}
		}

		if par.currentItem.typ == itemLeftBrace {
			if par.nextItem.typ < itemKeyWord {
				return par.reportError("identifier")
			}
		}

		if par.currentItem.typ == itemRightBrace {
			if par.nextItem.typ != itemEOF &&
				par.nextItem.typ != itemEntity {
				return par.reportError("entity")
			}
		}

		if par.currentItem.typ == itemString {
			if par.prevItem.typ == itemEntity {
				if par.nextItem.typ != itemLeftBrace {
					return par.reportError("{")
				}

				if par.currentItem.val == "#" {
					isCompiler = true
				} else {
					isCompiler = false
				}
				enityName = par.currentItem.val

			} else {
				if par.nextItem.typ != itemSemicolon {
					return par.reportError(";")
				}
				params = par.currentItem.val
			}

		}

		if par.currentItem.typ == itemEntity {
			if par.nextItem.typ != itemString {
				return par.reportError("entity name")
			}

			if par.currentItem.val == "#" {
				isCompiler = true
			} else {
				isCompiler = false
			}

			enityName = par.currentItem.val
		}

		if par.currentItem.typ > itemKeyWord {
			if par.nextItem.typ != itemEquals {
				return par.reportError("=")
			}

			identifier = par.currentItem.typ
		}
	}

	return nil
}

func newParser(file string) *parser {
	lex := newLexer(file)
	lex.analyze()
	par := parser{
		input:    lex.items,
		pos:      0,
		nextItem: lex.items[0],
	}
	return &par
}

func (par *parser) reportError(expected string) error {
	return errors.New("Syntax error on line " + string(par.currentItem.line) +
		": Expected " + expected + " , found " + par.nextItem.val)
}

var compilerDetails compiler
var fileDetails map[string]params

func fillCompilerDetails(identifier itemType, param string) {
	if identifier == itemBinary {
		compilerDetails.binary = param
	}
	if identifier == itemName {
		compilerDetails.name = param
	}
	if identifier == itemStart {
		compilerDetails.start = param
	}
	if identifier == itemLdFlags {
		compilerDetails.ldFlags = param
	}
	if identifier == itemIncludes {
		compilerDetails.includes = param
	}
	if identifier == itemOthers {
		compilerDetails.otherFlags = param
	}
}

func fillFileDetails(name string, identifier itemType, param string) {
	var temp params

	if identifier == itemFile {
		temp.file = param
	} else if param != "" {
		temp = fileDetails[name]
	}

	if param == "" {
		return
	}

	if identifier == itemDeps {
		paramArray := strings.Split(param, " ")
		temp.deps = paramArray
	}

	fileDetails[name] = temp
}

func init() {
	fileDetails = make(map[string]params)
}
