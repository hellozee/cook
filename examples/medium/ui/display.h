#pragma once

//Drawing the window using SDL

#include <string>
#include <iostream>
#include <SDL2/SDL.h>
#include <GL/glew.h>

class Display{

public:
    //Takes in the width , height and title of the window
    Display(int width, int height, const std::string &title);
    ~Display();

    //Checking whether the window is closed or not
    inline bool isClosed(){return _isClosed;}
    //Display what is drawn
    void swapBuffers();
    //Clearing the screen for further drawing
    void clear(float r,float g,float b,float a);

    void manageEvents(SDL_Event event);

private:
    SDL_Window *_window;
    SDL_GLContext _glContext;

    bool _isClosed = false;

};