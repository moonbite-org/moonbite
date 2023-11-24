# If Statements

```
if (bool) {
  // executes if true
}

if (bool1) {
  // executes if bool1 true
}else if(bool2) {
  // executes if bool2 true
}

if (bool) {
  // executes if true
}else {
  // executes if not
}
```

# Match Expressions

The `match` control flow is an expression that might resolve to a value. The resulting type is a piped type like `String | Int`. The value itself can be referenced with the self `.` operator.

```
match (value) {
  . == true {
    return "some string"
  }
  . instanceof String {
    return 3
  }
}
```

> `match` expressions are type checked in that you cannot compile above program if type of `value` is not akin to `Bool | String`. The resulting value from this expression will have the type `String | Int` and the order of matches matter.

# Loops

## Unipartite Loop

Basically a while loop. While a value is `true` it executes the body.

```
for (true) {
  // infinite loop
}
```

## Bipartite Loop

Iterates over an iterable value with labels.

```
for (key, element of iterable) {
  // in the case of lists, key is the index
}

for (,element of iterable) {
  // you can omit any of the values on the left hand side
}

for (, of iterable) {
  // ...
}
```

To iterate over a value, the value must implement the `Iterable<T>` interface, it looks like this:

```
type IteratorResult<T> {
  value T
  is_done Bool
}

trait Iterable<T> {
  fun next() IteratorResult<T>
}
```

## Tripartite Loop

A traditional three part for loop.

```
for (var index = 0; index < length; index++) {
  // ...
}

for (var index = length; index > 0; index--) {
  // ...
}
```
