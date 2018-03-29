package main

import (
	"errors"
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
	/*
		isCompiler := false
		enityName := ""
		identifier := itemNULL
		params := ""
	*/
	for par.next().typ != itemEOF {

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
			} else {
				if par.nextItem.typ != itemSemicolon {
					return par.reportError(";")
				}
			}

		}

		if par.currentItem.typ == itemEntity {
			if par.nextItem.typ != itemEquals {
				return par.reportError("=")
			}
		}

		if par.currentItem.typ == itemBinary {
			if par.nextItem.typ != itemEquals {
				par.reportError("=")
			}
		}

		if par.currentItem.typ == itemName {
			if par.nextItem.typ != itemEquals {
				return par.reportError("=")
			}
		}

		if par.currentItem.typ == itemStart {
			if par.nextItem.typ != itemEquals {
				return par.reportError("=")
			}
		}

		if par.currentItem.typ == itemLdFlags {
			if par.nextItem.typ != itemEquals {
				return par.reportError("=")
			}
		}

		if par.currentItem.typ == itemIncludes {
			if par.nextItem.typ != itemEquals {
				return par.reportError("=")
			}
		}

		if par.currentItem.typ == itemOthers {
			if par.nextItem.typ != itemEquals {
				return par.reportError("=")
			}
		}

		if par.currentItem.typ == itemFile {
			if par.nextItem.typ != itemEquals {
				return par.reportError("=")
			}
		}

		if par.currentItem.typ == itemDeps {
			if par.nextItem.typ != itemEquals {
				return par.reportError("=")
			}
		}
	}

	return nil

}

func newParser(file string) *parser {
	lex := newLexer(file)
	lex.analyze()
	par := parser{
		input:       lex.items,
		pos:         0,
		currentItem: lex.items[0],
		prevItem:    item{typ: itemNULL, pos: -1, val: "", line: -1},
		nextItem:    lex.items[1],
	}
	return &par
}

func (par *parser) reportError(expected string) error {
	return errors.New("Syntax error on line " + string(par.currentItem.line) +
		": Expected " + expected + " , found " + par.nextItem.val)
}

var compilerDetails compiler
var fileDetails map[string]params

func main() {
	temp := `entity = #{
		binary = g++;
		name = GLWindow;
		start = main;
		ldflags = -lSDL2 -lGLEW -lGL -lSOIL;
		includes = ;
		others = -Wall -Wextra;
	}
	
	entity = main{
		file = main.cpp;
		deps = camera display mesh shader;
	}
	
	entity = camera{
		file = ui/camera.cpp;
		deps = ;
	}
	
	entity = display{
		file = ui/display.cpp;
		deps = ;
	}
	
	entity = mesh{
		file = draw/mesh.cpp;
		deps = ;
	}
	
	entity = shader{
		file = draw/shader.cpp;
		deps = ;
	}`

	par := newParser(temp)
	par.parse()

}
