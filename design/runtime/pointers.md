Pointers are just references to variables and only there to pass values and mutate them effectively. A `Reference` is always typed with its value and adheres to the type system. Depending on the kind of a variable (ie `const`, `var`), a corresponding `Reference` exists and again, adheres to the mutability rules. The `ref` and `deref` builtin functions are use for referencing and dereferencing values.

```
fun main() {
  var Int x = 4 // x instanceof Int
  var p_x = ref(x) // p_x instanceof MutableReference<Int>
  var y = deref(p_x) // y instanceof Int
}

fun main() {
  const Int x = 4 // x instanceof Int
  x++ // this is invalid
  var p_x = ref(x) // p_x instanceof ImmutableReference<Int>
  var y = deref(p_x) // y instanceof Int
  y++ // this is valid
}

fun main() {
  var Int x = 4 // x instanceof Int
  var p_x = ref(x) // p_x instanceof MutableReference<Int>
  refset(p_x, /* accepts only Int */, 5)
  var y = deref(p_x) // y instanceof Int and y == 5
}

fun main() {
  const Int x = 4 // x instanceof Int
  x++ // this is invalid
  var p_x = ref(x) // p_x instanceof ImmutableReference<Int>
  refset(p_x, 5) // invalid, p_x is an immutable reference.
  var y = deref(p_x) // y instanceof Int
  y++ // this is valid
}
```

> In the last example, the program will not compile because the `refset` function only accepts `MutableReference<T>` type as its first argument so it will produce a compile time type error.

> You don't have to use the `refset` function for structs.

```
type User {
  first_name String
  last_name String
}

fun main() {
  var user = User{
    first_name: "John",
    last_name: "Doe",
  }

  update_user(ref(user))
}

fun update_user(user MutableReference<User>) {
  user.first_name = "Jane"
}
```

In the above example, the `update_user` must accept a `MutableReference<User>`, like refset, an update on a `ImmutableReference<T>` is not possible. So the below example will raise a compile time error.

```
fun update_user(user ImmutableReference<User>) {
  user.first_name = "Jane"
}
```

Reference types are as follows:

- `Reference<T>`: any reference to a value typed T
- `MutableReference<T>`: a mutable reference to a value typed T
- `ImmutableReference<T>`: an immutable reference to a value typed T