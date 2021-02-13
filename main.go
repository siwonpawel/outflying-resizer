package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/siwonpawel/outflying-resizer/resizer"
)

var (
	threads int
	folder  string
	scale   float64
	output  string
)

func main() {
	start := time.Now()
	signals := registerListenerForKeyBindings()
	parseArguments()

	r := resizer.New(threads, scale, folder, output)
	wg, cancelFunc, err := r.StartProcessingWithCancel()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	go func() {
		select {
		case <-signals:
			fmt.Println("Gracefull shitdown...")
			cancelFunc()
		}
	}()

	wg.Wait()
	fmt.Printf("Processing done in %0.3v s\n", time.Since(start).Seconds())
}

func registerListenerForKeyBindings() <-chan os.Signal {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT)

	return sigc
}

func parseArguments() {

	flag.IntVar(&threads, "t", runtime.NumCPU(), "number of threads")
	flag.StringVar(&folder, "f", "", "folder containing files")
	flag.Float64Var(&scale, "s", 100.00, "scale in percentage")
	flag.StringVar(&output, "o", "", "output folder of processed files, default output overwrites files in place")
	flag.Parse()

	fmt.Println(folder, "||", output)
	if folder == "" || output == "" {
		flag.PrintDefaults()
		os.Exit(0)
	}
}
