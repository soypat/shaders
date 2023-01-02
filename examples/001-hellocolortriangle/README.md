![Hello colorful triangle](https://user-images.githubusercontent.com/26156425/210190698-e4f439a0-115f-4796-9616-2bf4d2a7ff81.png)

## GLSL pitfalls

* Variables can be optimized out. Be sure the 
variable you seek in a call to `GetAttribLocation` is used
meaningfully in the program so it is not optimized out. This causes a -1 index to be returned.

* `size` attribute is limited to values between 1..4