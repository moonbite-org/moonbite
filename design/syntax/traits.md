A trait is a set of methods to be implemented by a type that implements that trait. Traits can be generic and the methods use the standard function notation without the `fun` keyword and the function body.

```
trait Greeter {
  greet(name String) String
}

type Greet implements [Greeter] {}
```

Traits can be implemented by structs with the following syntax

```
fun for Greet greet(name String) String {
  // ...
}
```

Traits can inherit other traits by mimicing them similar to how types implement traits.

```
trait Greeter {
  greet(name String) String
}

trait Repeater {
  repeat(words String) String
}

trait LobbyBoy mimics [Greeter, Repeater] {
  welcome(name String) String
}
```

Now a type that implements the `LobbyBoy` trait must implement `greet`, `repeat` and `welcome` methods.