package pararun

import (
	"fmt"
	"io"
	"sync"

	"github.com/fatih/color"
	multierror "github.com/hashicorp/go-multierror"
	parawri "github.com/mix3/go-parawri"
)

type Func struct {
	Name       string
	Func       func() error
	formatName string
	w          io.Writer
}

type Runner struct {
	DoneStr string
	FailStr string
	ProgStr string
	funcs   []*Func
	c       int
	p       *parawri.Parallel
	maxLen  int
}

func NewRunner(c int) *Runner {
	if c < 1 {
		c = 1
	}
	return &Runner{
		DoneStr: color.GreenString("done"),
		FailStr: color.RedString("fail"),
		ProgStr: "...",
		funcs:   []*Func{},
		c:       c,
		p:       parawri.NewParallelStdout(),
		maxLen:  0,
	}
}

func (r *Runner) AddFunc(f *Func) {
	f.w = r.p.NewAppendWriter()
	r.funcs = append(r.funcs, f)
	if r.maxLen < len(f.Name) {
		r.maxLen = len(f.Name)
	}
}

func (r *Runner) Run() error {
	for _, f := range r.funcs {
		f.formatName = fmt.Sprintf(fmt.Sprintf("%%-%ds", r.maxLen), f.Name)
		fmt.Fprint(f.w, f.formatName)
	}

	var (
		err   error
		errCh = make(chan error)
	)
	go func() {
		var (
			wg sync.WaitGroup
			q  = make(chan *Func)
		)
		for i := 0; i < r.c; i++ {
			wg.Add(1)
			go r.run(&wg, q, errCh)
		}
		for _, f := range r.funcs {
			q <- f
		}
		close(q)
		wg.Wait()
		close(errCh)
	}()

	for e := range errCh {
		err = multierror.Append(err, e)
	}

	return err
}

func (r *Runner) run(wg *sync.WaitGroup, q chan *Func, errCh chan error) {
	defer wg.Done()
	for {
		f, ok := <-q
		if !ok {
			break
		}
		fmt.Fprintf(f.w, " %s ", r.ProgStr)
		if err := f.Func(); err != nil {
			errCh <- fmt.Errorf("%s: %v", f.formatName, err)
			fmt.Fprint(f.w, r.FailStr)
		} else {
			fmt.Fprint(f.w, r.DoneStr)
		}
	}
}
