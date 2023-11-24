Saturables are the heart of the `dynamic saturization` concept. They are a generic trait that represents the type itself with a single additional method.

```
trait Saturable<T> {
  fun default() T
}

type Int implements [Saturable<int>] int

fun for Int default() int {
  return raw(0)
}

type Bool implements [Saturable<bool>] bool

fun for Bool default() bool {
  return raw(false)
}

// ...

fun raw<T string | bool | number | rune>(value Saturable<T>) T
```