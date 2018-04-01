#pragma once

#include <glm/glm.hpp>
#include <GL/glew.h>
#include <SOIL/SOIL.h>
#include <glm/glm.hpp>
#include <glm/gtc/matrix_transform.hpp>
#include <glm/gtc/type_ptr.hpp>
#include <vector>


class Mesh{

public:
    Mesh(std::vector<GLfloat> meshData);
    ~Mesh();

    void Draw(GLuint program);
    void addTexture(const char * path);
    void Translate(glm::vec3 transVec);
    void Rotate(float angle, glm::vec3 axis);
    void Scale(glm::vec3 scaleVec);

private:

    GLuint _vbo,_vao,_texture;

    glm::mat4 _model;

    unsigned int _drawCount;
};