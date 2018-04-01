#include "shader.h"

Shader::Shader(const std::string &vertexShaderSource,const std::string &fragmentShaderSource)
{
    program = glCreateProgram();
    _shaders[0] = CreateShader(LoadShader(vertexShaderSource),GL_VERTEX_SHADER);
    _shaders[1] = CreateShader(LoadShader(fragmentShaderSource),GL_FRAGMENT_SHADER);

    for(unsigned int i=0;i<2;i++){
        glAttachShader(program,_shaders[i]);
    }

    glLinkProgram(program);

    CheckShaderError(program,GL_LINK_STATUS,true,"Error Shader program failed to link");

    glValidateProgram(program);

    CheckShaderError(program,GL_VALIDATE_STATUS,true,"Error Shader program is invalid");
}

Shader::~Shader()
{
    for(int i=0;i<2;i++){
        glDetachShader(program,_shaders[i]);
        glDeleteShader(_shaders[i]);
    }
    glDeleteProgram(program);
}

void Shader::Bind()
{
    glUseProgram(program);
}

GLuint Shader::CreateShader(const std::string &text, GLenum shaderType)
{
    GLuint shader = glCreateShader(shaderType);

    if(shader == 0){
        std::cerr << "Error : Shader Creation failed." << std::endl;
    }

    const GLchar* shaderSource;
    GLint shaderLength;

    shaderSource = text.c_str();
    shaderLength = text.length();

    glShaderSource(shader,1,&shaderSource,&shaderLength);
    glCompileShader(shader);

    CheckShaderError(shader,GL_COMPILE_STATUS,false,"Error Shader compilation failed");

    return shader;
}

std::string Shader::LoadShader(const std::string& fileName)
{
    std::ifstream file;
    file.open((fileName).c_str());

    std::string output;
    std::string line;

    if(file.is_open())
    {
        while(file.good())
        {
            getline(file, line);
			output.append(line + "\n");
        }
    }
    else
    {
		std::cerr << "Unable to load shader: " << fileName << std::endl;
    }

    return output;
}

void Shader::CheckShaderError(GLuint shader, GLuint flag, bool isProgram, const std::string& errorMessage)
{
    GLint success = 0;
    GLchar error[1024] = { 0 };

    if(isProgram)
        glGetProgramiv(shader, flag, &success);
    else
        glGetShaderiv(shader, flag, &success);

    if(success == GL_FALSE)
    {
        if(isProgram)
            glGetProgramInfoLog(shader, sizeof(error), NULL, error);
        else
            glGetShaderInfoLog(shader, sizeof(error), NULL, error);

        std::cerr << errorMessage << ": '" << error << "'" << std::endl;
    }
}
