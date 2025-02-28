# A Lua VM pool in Go

A Lua VM pool for [gopher-lua](https://github.com/yuin/gopher-lua).


## Example:

```go
package main 

import (
    "log"
    lpool "github.com/epikur-io/glua-pool"
)

func main() {
    pool := lpool.NewPool(10, nil)

    // get a VM:
    luaVM := pool.Acquire()

    // do stuff...

    // release VM
    pool.Release(luaVM)

    // get a VM or timeout after 1 second:
    luaVM, err := pool.AcquireWithTimeout(time.Seconds * 1)
    if (err != nil) {
        log.Println("error:", err)
    }else{
        // do stuff...

        // release VM
        pool.Release(luaVM)
    }
}
```