package main

import (
	"fmt"
	"io"

	"time"

	"github.com/rpoletaev/respio"
)

func main() {
	r, w := io.Pipe()

	rw := respio.NewWriter(w)
	rr := respio.NewReader(r)

	go func() {
		for {
			time.Sleep(5 * time.Second)
			cmd, prs, err := rr.ReadCommand()
			if err != nil {
				println(err.Error())
			}

			println("command name: ", cmd)
			fmt.Printf("parameters %v\n", prs)
		}
	}()

	rw.SendCmd("raz", []interface{}{"dva"})
	rw.Flush()

	time.Sleep(1 * time.Second)
}
