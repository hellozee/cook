#pragma once

#include <GL/glew.h>
#include <glm/glm.hpp>
#include <glm/gtc/matrix_transform.hpp>
#include <glm/gtc/type_ptr.hpp>
#include <SDL2/SDL.h>
#include <iostream>
#include <vector>


class Camera{

public:
	Camera(glm::vec3 cameraPos, glm::vec3 cameraFront, glm::vec3 cameraUp, GLfloat fov, GLfloat ratio, GLfloat nearPlane, GLfloat farPlane);
	~Camera();

	void Use(GLuint program);
	void manageEvents(SDL_Event event);
	void doMovement(unsigned int currentFrame);

	glm::vec3 cameraPos;

private:
	glm::mat4 _projection, _view;
	glm::vec3 _cameraFront, _cameraUp;

	unsigned int _deltaTime, _lastFrame;

	GLfloat _pitch, _yaw, _sensitivity;
	GLfloat _xoffset,_yoffset;

	bool _firstMouse,_out;
	std::vector<bool>_keys;

	GLfloat _cameraSpeed;

	void changeView(SDL_Event event);

};