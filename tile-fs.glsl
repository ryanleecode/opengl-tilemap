#version 330 core
out vec4 outputColor;

in vec2 fragTexCoord;
uniform sampler2D texture1;

void main() {
 // outputColor = vec4(fragTexCoord, 0.0f, 1.0f);
   outputColor =  texture(texture1, fragTexCoord);
}
