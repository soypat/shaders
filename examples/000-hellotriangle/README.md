# VBOs
Vertex attributes are views into the VBO, which
may be much larger than the data that is interesting to the renderer.

# `glVertexAttribPointer` behaviour
This function is quite counter intuitive, even when one has
a good grasp on pointers, structs and offsets.

![glVertexAttribPointer explanation](https://i.stack.imgur.com/jh89v.png)

TODO: I still don't understand why the heck I'm calling

```go
gl.VertexAttribPointerWithOffset(vertAttrib, 2, gl.FLOAT, false, 2*attrSize, 0)
```
With size=2... Why?
My best guess is that each "vertex" is composed of two `gl.FLOAT`'s, so the size of the vertex is 2? 
This may be related to the `DrawArrays` call which takes a `count` parameter. Which is 3:

```go
// Inside render loop...
gl.DrawArrays(gl.TRIANGLES, 0, 3)
```
So maybe DrawArrays iterates over the array with
the following loop condition(?):

```go
iLimit := size*count

for i:=first; i < iLimit; i++ {
    offset := i*stride + first
    // ...
}
```



