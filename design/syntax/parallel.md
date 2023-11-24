Lng has Go style coroutines that compute given processes in parallel. The communication between these coroutines can be done via a `Channel`.

```
fun main() {
  corout {
    console.log("This will run in parallel")
  }

  console.log("This will run immediately")
}
```

> Although the last `console.log` call will run immediately, there is no guarantee that it will run first. It only means that the `corout` will not run sequentially.

You can block the main thread by expecting a value from a channel.

```
use time

fun main() {
  var Channel<Int> data

  corout {
    time.sleep(time.second * 5)
    data <- 1
  }

  console.log("Start")
  var count = <- data
  console.log("End")
}
```

This program will log "Start", wait for 5 seconds and then log "End".

> The `Channel` is a runtime construct that is provided.

> Other specifics to be determined...