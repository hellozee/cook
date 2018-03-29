package parser

import (
	"errors"
	"strings"
)

type compiler struct {
	Binary     string
	Name       string
	Start      string
	LdFlags    string
	Includes   string
	OtherFlags string
}

type params struct {
	File string
	Deps []string
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

func (par *parser) Parse() error {

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

func (par *parser) reportError(expected string) error {
	return errors.New("Syntax error on line " + string(par.currentItem.line) +
		": Expected " + expected + " , found " + par.nextItem.val)
}

func NewParser(file string) parser {
	lex := newLexer(file)
	lex.analyze()
	par := parser{
		input:    lex.items,
		pos:      0,
		nextItem: lex.items[0],
	}
	return par
}

//CompilerDetails ...
var CompilerDetails compiler

//FileDetails ...
var FileDetails map[string]params

func fillCompilerDetails(identifier itemType, param string) {
	if identifier == itemBinary {
		CompilerDetails.Binary = param
	}
	if identifier == itemName {
		CompilerDetails.Name = param
	}
	if identifier == itemStart {
		CompilerDetails.Start = param
	}
	if identifier == itemLdFlags {
		CompilerDetails.LdFlags = param
	}
	if identifier == itemIncludes {
		CompilerDetails.Includes = param
	}
	if identifier == itemOthers {
		CompilerDetails.OtherFlags = param
	}
}

func fillFileDetails(name string, identifier itemType, param string) {
	var temp params

	if identifier == itemFile {
		temp.File = param
	} else if param != "" {
		temp = FileDetails[name]
	}

	if param == "" {
		return
	}

	if identifier == itemDeps {
		paramArray := strings.Split(param, " ")
		temp.Deps = paramArray
	}

	FileDetails[name] = temp
}

func init() {
	FileDetails = make(map[string]params)
}
