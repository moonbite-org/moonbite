# Cardinal & Main Types

There are few cardinal types that make up the lng. These are:

- string
- bool
- rune
- int - int8 - int16 - int32 - int64
- uint - uint8 - uint16 - uint32 - uint64
- float32 - float64

From these, the main types are derived with which you can use a range of methods. They are expanded to be used throughout you program. These are mainly:

- String
- PathBuff
- Bool
- Rune
- Int - Int8 - Int16 - Int32 - Int64
- Uint - Uint8 - Uint16 - Uint32 - Uint64
- Byte
- Float32 - Float64

# Creating Types

You can create your own types with the `type` keyword.

```
type TokenKind Int
```

> You can use the `int` cardinal to represent your `TokenKind` type as well but you probably do not want to use cardinal types unless you have to in your program.

> The main reason for the above note is becuase you would have to implement a lot of traits to use the cardinal types on their own. For instance, the `console.log` function will not accept any cardinal types because they do not implement the `Printable` trait.

## Structs

You can create structs with `{ ... }` syntax and the `type` keyword.

```
type TimeLord {
  name String
  age Uint
}
```

An instance of this type can be instantiated like this:

```
TimeLord{name: "The Doctor", age: 903}
```

## Type Literals

You can create type literals to use them in enum-like ways. A type literal is created with the parent type name. When created, this becomes its own type that can represent both its own type and its parent type.

```
type UserKind String("admin") | String("user")
```

The `|` (pipe) operator denotes that the type `UserKind` can be either of these values. You can further use this type like this:

```
type User {
  first_name String
  last_name String
  age Uint
  kind UserKind
}
```

> The `UserKind` type can replace a `String` but not any `String` can replace the `UserKind` type. For instance, a function that accepts a `String` will accept a `UserKind` but not vice versa.

## Generic Types

Types can be generic, shaping itself with a given type. Generics are denoted with the `< ... >` syntax and there may be multiple generics seperated with commas. Generics can have constraints.

```
type Slice<T> {
  value T
}

type Deferred<T, W Warning> {
  value T
  warning W
}
```

# Type operators

Type operators are operators that when applied to one or more type, result in another type.

## Pipe `|` Operator

The pipe `|` operator works like a logical or operator, picking either of the types. It is binary operator that takes two types.

## Join `&` Operator
The join `&` operator conjoins two structs. If provided types are not structs, the program will error and not compile. If there are conflicting fields in the structs, the last struct provided will take precedence.