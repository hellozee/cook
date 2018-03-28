package main

import "fmt"

func main() {
	temp := `entity ${
		binary = g++;
		name = GLWindow;
		start = main;
		ldflags = -lSDL2 -lGLEW -lGL -lSOIL;
		includes = ;
		others = -Wall -Wextra;
	}
	
	entity main{
		file = main.cpp;
		deps = camera display mesh shader;
	}
	
	entity camera{
		file = ui/camera.cpp;
		deps = ;
	}
	
	entity display{
		file = ui/display.cpp;
		deps = ;
	}
	
	entity mesh{
		file = draw/mesh.cpp;
		deps = ;
	}
	
	entity shader{
		file = draw/shader.cpp;
		deps = ;
	}`
	l := lexer{
		input: temp,
		pos:   0,
		start: 0,
		width: len(temp),
		line:  0,
	}

	l.analyze()

	fmt.Println(l.items)
}
