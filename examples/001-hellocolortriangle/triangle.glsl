#shader fragment
#version 330

// v_vert_color is received from the
// vertex shader. That's how the 
// GPU pipeline works. Frags follow the vertex.
in vec4 v_vert_color;

out vec4 color;

void main() {
	color = vec4(v_vert_color.r, v_vert_color.g, v_vert_color.b, 1.0);
	color = color+ vec4(1.0, 0.0 ,0.0 ,1.0);
}

#shader vertex
#version 330

in vec3 vert;
in vec4 vert_color;
out vec4 v_vert_color;

void main() {
	// v_vert_color is piped
	// to the fragment shader.
	v_vert_color = vec4(vert_color.r, vert_color.g, vert_color.b, 1.0);
	// Declare positions of our vertices.
	gl_Position =  vec4(vert, 1.0);
}