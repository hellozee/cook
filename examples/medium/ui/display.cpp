#include "display.h"

Display::Display(int width, int height, const std::string &title)
{
    //Initializing SDL
    SDL_Init(SDL_INIT_EVERYTHING);

    //Setting up OpenGL Attributes
    SDL_GL_SetAttribute(SDL_GL_RED_SIZE,8);
    SDL_GL_SetAttribute(SDL_GL_BLUE_SIZE,8);
    SDL_GL_SetAttribute(SDL_GL_GREEN_SIZE,8);
    SDL_GL_SetAttribute(SDL_GL_ALPHA_SIZE,8);
    SDL_GL_SetAttribute(SDL_GL_BUFFER_SIZE,32);
    SDL_GL_SetAttribute(SDL_GL_DOUBLEBUFFER,1);

    SDL_GL_SetAttribute(SDL_GL_CONTEXT_MAJOR_VERSION, 3);
    SDL_GL_SetAttribute(SDL_GL_CONTEXT_MINOR_VERSION, 3);

    _window = SDL_CreateWindow(title.c_str(),SDL_WINDOWPOS_CENTERED,SDL_WINDOWPOS_CENTERED,width,height,SDL_WINDOW_OPENGL);
    _glContext = SDL_GL_CreateContext(_window);

    SDL_SetRelativeMouseMode(SDL_TRUE);
    SDL_WarpMouseInWindow(_window, width/2, height/2);

    glewExperimental = GL_TRUE; //oh boy you gave me a lot of headache
    GLenum status = glewInit();

    if(status != GLEW_OK){
        std::cerr << "GLEW failed to initialize." << std::endl;
    }

    glEnable(GL_DEPTH_TEST);
}

Display::~Display()
{
    SDL_GL_DeleteContext(_glContext);
    SDL_DestroyWindow(_window);
    SDL_Quit();
}

void Display::clear(float r,float g,float b,float a)
{
    glClearColor(r,g,b,a);
    glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT);
}

void Display::swapBuffers()
{
    SDL_GL_SwapWindow(_window);
}

void Display::manageEvents(SDL_Event event)
{
    if(event.type == SDL_QUIT){
        _isClosed = true;
    }

    if(event.type == SDL_KEYDOWN){
        switch (event.key.keysym.sym){
            case SDLK_ESCAPE:
                _isClosed = true;
                break;
            default:
                break;
        }
    }
}