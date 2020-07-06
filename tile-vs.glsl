#version 330 core

layout (location = 0) in vec2 position;
layout (location = 1) in vec2 vertTexCoord;

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

out vec2 fragTexCoord;

void main() {
  gl_Position = projection * camera * model * vec4(position, 0., 1.);
  fragTexCoord = vertTexCoord;
}
