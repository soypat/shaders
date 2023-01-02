#

## GLSL pitfalls

* Variables can be optimized out. Be sure the 
variable you seek in a call to `GetAttribLocation` is used
meaningfully in the program so it is not optimized out. This causes a -1 index to be returned.