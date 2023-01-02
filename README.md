# shaders
I'm learning shaders, why don't you too?

### Hello Magenta World
The following "shader" is a program that sets the whole viewport
to magenta. So you get a magenta screen.
```glsl
void main() {
    gl_FragColor = vec4(pi/4.,0.0,1.0,1.0);
}
```
It is a **fragment shader**. This means it is executed for every fragment (usually a single pixel) 
in the data pipeline (which may or may not be all of the screens pixels).
`gl_FragColor` is an output variable to the function. We set it to the pixel's
desired color. We may get the pixel's coordinates on screen with [`gl_FragCoord`](https://registry.khronos.org/OpenGL-Refpages/gl4/html/gl_FragCoord.xhtml) (`vec4`)

Below is the same program but with a bunch of other GLSL
language constructs to get familiar:
```glsl
// Data detailed here is for OpenGL.
// This is an example of a fragment shader.
// It sets a certain pixel to a solid color via the
// pixel-specific variable: gl_FragColor

#ifdef GL_ES
// We can choose the global precision of the shader
// Also, don't end comments before precision call with "precision".
precision lowp float;
#endif

// Uniforms are defined at the top of shader
// After assigning the default floating poing precision
// Uniforms are inputs which are equal for all
// threads and necessarily set to read only.
// Time in seconds since shader loaded.
uniform float u_time;
// Mouse position in screen pixels.
uniform vec2 u_mouse;

// Canvas size (width, height).
uniform vec2 u_resolution;

// These uniform variables are the same as the above
// but are integers (????).

// Time in seconds since load
// Viewport resolution in pixels.
uniform vec3 iResolution; 

// Mouse pixel coordinates. xy:current. zw:click
uniform vec4 iMouse;

// Shader playback time in seconds.
uniform float iTime;

// You can define macros. 
// They need not be limited to numbers. 
// Notice this line has no semicolon!
// The define label is `pi` and everything that follows
// after the space is what is replaced whenever
// the label is encountered in the code.
#define pi 3.1415926535

// You also may not initialize uniform variables:
// uniform float v=1.0; // ERROR: cannot initialize this type of qualifier

// GPU available hardware accelerated functions.
// sin(x), cos(x), tan(x), asin(x), 
// acos(x), atan(x, [y]), pow(x,y), exp(x), 
// log(x), sqrt(), abs(), sign(), 
// floor(), ceil(), fract(), mod(), 
// min(), max() and clamp().
// rand(x)

void main() {
    gl_FragColor = vec4(pi/4.,0.0,1.0,1.0);
    // Code below: EPILEPSY WARNING
    // float frequency = u_mouse.y/60.0;
    // float t = u_time;
    // gl_FragColor = vec4(abs(sin(pi*frequency*t)),0.0,0.0,1.0);
}
```