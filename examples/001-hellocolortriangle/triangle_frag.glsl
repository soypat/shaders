#version 330

in vec4 v_vert_color;

out vec4 color;

void main() {
	color = vec4(v_vert_color.r, v_vert_color.g, v_vert_color.b, 1.0);
	color = color+ vec4(1.0, 0.0 ,0.0 ,1.0);
}