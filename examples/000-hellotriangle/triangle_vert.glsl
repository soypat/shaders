#version 330

in vec3 vert;

void main() {
	gl_Position = vec4(vert.xyz, 1.0);
}