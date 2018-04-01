#version 330 core

in vec3 fragPos;
in vec3 Normals;

out vec4 color;

uniform vec3 lightColor;
uniform vec3 cubeColor;
uniform vec3 lightPos;
uniform vec3 viewPos;

void main(){

    float ambientStrength = 0.1f;
    float specularStrength = 1.0f;

    vec3 ambient = ambientStrength * lightColor;

    vec3 norm = normalize(Normals);
    vec3 lightDir = normalize(lightPos - fragPos);
    vec3 viewDir = normalize(viewPos - fragPos);
    vec3 reflectDir = reflect(-lightDir,norm);

    float diff = max(dot(norm, lightDir), 0.0);
    float spec = pow(max(dot(viewDir,reflectDir),0.0),256);

    vec3 specular = specularStrength * spec * lightColor;
  	vec3 diffuse = diff * lightColor;
  	vec3 result = (ambient + diffuse + specular) * cubeColor;

    color = vec4(result,1.0f);
}