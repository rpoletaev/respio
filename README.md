# respio
allow read and write RESP data http://redis.io/topics/protocol
```Go
package main

import(
"fmt"
"time"
"github.com/rpoletaev/respio"
)

func main() {
	r, w := io.Pipe()

	rw := respio.NewWriter(w)
	rr := respio.NewReader(r)

	go func() {
		for {
			cmd, prs, err := rr.ReadCommand()
			if err != nil {
				println(err.Error())
			}

			println("command name: ", cmd)
			fmt.Printf("parameters %v\n", prs)
		}
	}()

	rw.SendCmd("raz", []interface{}{"dva"}) //send command in RESP format like: *2\r\n$3\r\n114 97 12\r\n$3\r\n100 118 97\r\n
	rw.Flush()

	time.Sleep(1 * time.Second)
}

```
