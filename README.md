## Console input logger for Go (for Windows)

---

### About coninlogger

- Use coninlogger to log console input events (for Windows).

---

### Example

- This program listens for console input events for 3 seconds.

```go
package main

import (
	"fmt"
	"github.com/jeet-parekh/coninlogger"
	"time"
)

func main() {
	inl := coninlogger.NewConsoleInputLogger(4)
	inlmsg := inl.GetMessageChannel()
    inl.Start()
    
	go func() {
		time.Sleep(time.Second * 3)
		inl.Stop()
    }()
    
	for v := range inlmsg {
		fmt.Printf("%+v\n", v)
	}
}
```

---
