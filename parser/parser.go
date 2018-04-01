package parser

import (
	"errors"
	"strconv"
	"strings"
)

//compiler Data Structure to hold the details for the compiler
type compiler struct {
	Binary     string
	Name       string
	Start      string
	LdFlags    string
	Includes   string
	OtherFlags string
}

//params  Data Structure to hold the details for the file entity
type params struct {
	File string
	Deps []string
}

//Parser Data Structure for holding parser details and parsed object
type Parser struct {
	input           []item
	pos             int
	currentItem     item
	prevItem        item
	nextItem        item
	CompilerDetails compiler
	FileDetails     map[string]params
}

//next Function to shift to the next item in the list
func (par *Parser) next() item {
	par.prevItem = par.currentItem
	par.currentItem = par.nextItem
	par.pos++
	par.nextItem = par.input[par.pos]
	return par.currentItem
}

//Parse Function to parse the the Recipe File
func (par *Parser) Parse() error {

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
			}
		}

		if par.currentItem.typ == itemSemicolon {
			if par.nextItem.typ < itemKeyWord &&
				par.nextItem.typ != itemRightBrace {
				return par.reportError("}")
			}

			if isCompiler {
				par.fillCompilerDetails(identifier, params)

			} else {
				par.fillFileDetails(enityName, identifier, params)
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

//reportError Function for reporting syntax errors
func (par *Parser) reportError(expected string) error {
	return errors.New("Syntax error on line " + strconv.Itoa(par.currentItem.line) +
		": Expected " + expected + " , found " + par.nextItem.val)
}

//NewParser Helper function to create a parser
func NewParser(file string) Parser {
	lex := newLexer(file)
	lex.analyze()
	par := Parser{
		input:       lex.items,
		pos:         0,
		nextItem:    lex.items[0],
		FileDetails: make(map[string]params),
	}
	return par
}

//fillCompilerDetails Function to store the compiler details
func (par *Parser) fillCompilerDetails(identifier itemType, param string) {
	if identifier == itemBinary {
		par.CompilerDetails.Binary = param
	}
	if identifier == itemName {
		par.CompilerDetails.Name = param
	}
	if identifier == itemStart {
		par.CompilerDetails.Start = param
	}
	if identifier == itemLdFlags {
		par.CompilerDetails.LdFlags = param
	}
	if identifier == itemIncludes {
		par.CompilerDetails.Includes = param
	}
	if identifier == itemOthers {
		par.CompilerDetails.OtherFlags = param
	}
}

//fillFileDetails  Function to fill the file details
func (par *Parser) fillFileDetails(name string, identifier itemType, param string) {
	var temp params

	if identifier == itemFile {
		temp.File = param
	} else if param != "" {
		temp = par.FileDetails[name]
	}

	if param == "" {
		return
	}

	if identifier == itemDeps {
		paramArray := strings.Split(param, " ")
		temp.Deps = paramArray
	}

	par.FileDetails[name] = temp
}
