Functions are declared with a `fun` keyword. The function paratemers and the return type are always present and typed unless the function does not return or take anything.

```
fun greet(name string) string {
  return format("Hello %s", name)
}

fun say_hello(name string) {
  console.log("Hello %s", name)
}
```