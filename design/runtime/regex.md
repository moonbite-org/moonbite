```
package main

use regexp

fun main() {
  var identifier = regexp.compile("[a-zA-Z][a-zA-Z0-9]")
  identifier.match("0test")
  identifier.match("test0")
}
```