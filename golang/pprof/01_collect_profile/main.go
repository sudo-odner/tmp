package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"sync"
	"sync/atomic"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	f, err := os.Create(filepath.Join(wd, "cpu.out"))
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatalln(err)
	}
	defer pprof.StopCPUProfile()

	wg := new(sync.WaitGroup)
	for range runtime.NumCPU() / 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for range int(10e6) {
				root()
			}
		}()
	}

	wg.Wait()
}

var v atomic.Int32

//go:noinline
func root() {
	f2()
	f3()
	f4()
}

//go:noinline
func f2() {
	v.Add(1)
}

//go:noinline
func f3() {
	v.Add(1)
}

//go:noinline
func f4() {
	v.Add(1)
}
