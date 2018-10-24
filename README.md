# Cook 

A small and fast build system for your C and C++ projects. 

### <u>Installation</u>

```bash
git clone https://github.com/hellozee/cook
cd cook
make
```



### <u>How to use it?</u>

It is simple, just like a we have `Makefile` for make, here we will use a `Recipe` file for doing the same thing. 

#### <u>Syntax:</u>

Cook is pretty strict on syntax, we have little blocks of what I call entity, there are 2 kinds of entity, one is normal entity and another is a special entity for the compiler. Entities have have properties which are called identifiers, like the normal entities have `file` and `deps` identifiers.

- `file ` : Location of the file which this entity is about
- `deps` : Dependencies of the above mentioned file

 On the other hand the special entity for compiler has more than just 2 identifiers, which are `binary`  `name` `start`  `ldflags` `includes`  `others`.

- `binary` : The compiler binary which will be used to compile the program
- `name` : Name of the generated executable
- `start` :  The root entity, or better say the main entity 
- `ldflags` : Linker directives, enter them as you would enter them in the terminal
- `includes` : Same as Linker directives but for mentioning special include folders.
- `others` : Some other flags like `-Wall`

So, how do we separate normal entity from the special compiler entity? Simple, there is special name given to this entity which is `#` , it is reserved and hence can't be used for normal entities, also naming an entity is compulsory as it also be used for figuring out the dependencies since the same name would be used to in the `deps` identifier. Linebreaks are used as delimiter but you can use `;` if you want to mention multiple identifiers in a single line. 

##### Example:

```
entity #{
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
}
```

This is the `Recipe` file I am using for practicing `OpenGL` , you can try it on my OpenGL Practice [repository](https://github.com/hellozee/gl-practice) also this is available in examples folder.

After creating the `Recipe` file you just execute Cook in the directory containing the `Recipe` file. 

**Note:**  This only works for C/C++

### Usage Flags:

- `--help` : To show this help message
- `--clean` : To clean the cached data

### ToDos:

- Implement `+=` and `-=` operator
- Add an `import` feature to import other Recipes and recursively cook them
