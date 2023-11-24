Lng does not have runtime errors, only recoverable warnings. This is possible becuase of a concept called `dynamic saturation`. The value of a given variable is always set to a default value and its saturated state is observed. When a function or any procedure were to develop an unexpected result, the default value for the function return type is produced for the function call expression and a warning is raised. You can further recover from these warnings by warning the parent stack or quitting the runtime entirely, basically bubbling any warnings. Alternatively, you can take action based on the warning with `if`, `switch` or `match` statements.

# Raising a Warning

Let's assume you have a function that has to accept an `Int8` value that you would like to treat as a `Uint8` making it a non negative nibble. You would raise a warning if a negative value is passed.

```
fun nibble(value Int8) {
  if (value < 0) {
    warn("Cannot accept a negative value".(Warning))
  }
}
```

You would invoke the `warn` function with a `Warning` string (A string that is casted to `Warning`). This would immediately stop the execution of a function like a `return` or a `yield` statement and bubble the warning up.

# Recovering from a Warning

Let's try to recover from above warning by using the `or` keyword.

```
use os

fun main() {
  nibble(-1) or os.exit(1)
}
```

This will stop the execution process if the nibble function raises a warning, which it will. You can also react to the warning with a `match` statement.

# Caret `^` Operator

The caret `^` special operator provides a simple shorthand for the raised warning that you might need to handle. You can treat it like a normal variable, query its type, read its value or bubble it to upper stack.

```
use os

// bubble
fun another_fun() {
  nibble(-1) or warn(^) 
}

// react to it
fun main() {
  nibble(-1) or console.log(^)
}

// or
fun main() {
  os.read_file("non-existant.file") or match (^) {
    . instanceof os.NotFoundWarning -> {
      console.log("not found")
    }
    . instanceof os.ReadWarning -> {
      console.log("file is corrupt")
    }
    default -> {
      console.log("some other warning")
    }
  }
}
```

When a warning is raised, the default value is returned from the function, if a return value exists at all. You can return other values with a `match` statement as well.

> The fact that the runtime will not produce an error or stop its execution becuase of lack of errors should not mean lng programs are safer. It means it is more resillient to faults for its authors and the resultings programs' users. The linter can be configured to show uncaught warnings a function may raise and you should always handle them. A function's warning type is inferred and enforced, you cannot check if the `os.read_file` function's raised warning is an instance of some `Arbitrary` warning, it will not compile.