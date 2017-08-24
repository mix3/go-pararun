package main

import (
	"fmt"
	"os"
	"time"

	pararun "github.com/mix3/go-pararun"
)

func main() {
	r := pararun.NewRunner(2)
	//r.DoneStr = color.GreenString("ok")
	//r.FailStr = color.RedString("ng")
	//r.ProgStr = ">>>"
	for _, f := range []*pararun.Func{
		&pararun.Func{
			Name: "e1",
			Func: func() error {
				time.Sleep(time.Second * 1)
				return nil
			},
		},
		&pararun.Func{
			Name: "e2",
			Func: func() error {
				time.Sleep(time.Second * 3)
				return fmt.Errorf("ERROR!")
			},
		},
		&pararun.Func{
			Name: "e3",
			Func: func() error {
				time.Sleep(time.Second * 2)
				return nil
			},
		},
		&pararun.Func{
			Name: "e4",
			Func: func() error {
				time.Sleep(time.Second * 3)
				return nil
			},
		},
		&pararun.Func{
			Name: "e5",
			Func: func() error {
				time.Sleep(time.Second * 1)
				return fmt.Errorf("ERROR!")
			},
		},
		&pararun.Func{
			Name: "e6",
			Func: func() error {
				time.Sleep(time.Second * 2)
				return nil
			},
		},
	} {
		r.AddFunc(f)
	}
	fmt.Fprintln(os.Stderr, "\n", r.Run())
}
