//Includes
#include "ui/display.h"  //For drawing the window, uses SDL
#include "draw/mesh.h"   //For generating the meshes
#include "draw/shader.h" //Compiling Shaders
#include "ui/camera.h"   //Camera Controls
#include "mesh/cube.h"   //The default cube

int main(){

    const int WIDTH = 800, HEIGHT = 600;

    GLfloat ratio = (GLfloat) WIDTH/ (GLfloat) HEIGHT;

    Display disp(WIDTH,HEIGHT,"Testing"); // SDL Window for drawing things

    //Intializing mesh
    Mesh cube(defaultCube);
    //cube.addTexture("textures/container.jpg"); // Not required for now

    cube.Rotate(45.0f,glm::vec3(0.0f,1.0f,0.0f));
    cube.Translate(glm::vec3(-1.0f,0.0f,-1.0f));

    Mesh light(defaultCube);

    glm::vec3 lightPos(1.2f, 1.0f, 2.0f);

    light.Scale(glm::vec3(0.2f));
    light.Translate(lightPos);

    //Fetching Shaders
    Shader shader("shaders/vertexCube.glsl","shaders/fragmentCube.glsl");
    Shader lightShader("shaders/vertexLight.glsl","shaders/fragmentLight.glsl");

    //Creating the camera
    Camera cam(glm::vec3(0.0f, 0.0f, 6.0f), glm::vec3(0.0f, 0.0f, -1.0f), glm::vec3(0.0f, 1.0f, 0.0f), 45.0f, ratio, 0.1f, 100.0f);

    while(!disp.isClosed()){
        //Clearing the screen to draw
        disp.clear(0.2f,0.2f,0.2f,1.0f);

        SDL_Event event;

        if(SDL_PollEvent(&event)){
            disp.manageEvents(event);
            cam.manageEvents(event);
        }

        //Attaching the shader to the program
        shader.Bind();
        cam.Use(shader.program);
        //Drawing the mesh

        glUniform3f(glGetUniformLocation(shader.program,"cubeColor"),1.0f, 0.5f, 0.31f);
        glUniform3f(glGetUniformLocation(shader.program,"lightColor"),1.0f,1.0f,1.0f);
        glUniform3f(glGetUniformLocation(shader.program,"lightPos"),lightPos.x,lightPos.y,lightPos.z);
        glUniform3f(glGetUniformLocation(shader.program,"viewPos"), cam.cameraPos.x, cam.cameraPos.y, cam.cameraPos.z);

        cube.Draw(shader.program);

        lightShader.Bind();
        cam.Use(lightShader.program);
        light.Draw(lightShader.program);

        cam.doMovement(SDL_GetTicks());

        //Now Displying what is drawn
        disp.swapBuffers();
    }

    return 0;
}
