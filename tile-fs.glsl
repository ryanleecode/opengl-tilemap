#version 410 core
out vec4 outputColor;

in vec2 fragTexCoord;
uniform sampler2D tileAtlas;

void main() {
   outputColor =  texture(tileAtlas, fragTexCoord);
}
