#version 330

in vec3 vert;
in vec4 vert_color;
out vec4 v_vert_color;

void main() {
	v_vert_color = vec4(vert_color.r, vert_color.g, vert_color.b, 1.0);
	gl_Position =  vec4(vert, 1.0);
}