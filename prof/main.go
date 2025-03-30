package main

import (
	"context"
	"flag"
	"io"
	"os"
	"runtime/pprof"

	"github.com/localhots/riff"
	"github.com/localhots/riff/ctx/log"
)

func main() {
	cpuprofile := flag.String("cpuprofile", "", "Write cpu profile to file")
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	log.Setup(riff.Config{
		Level:           riff.LevelDebug,
		Output:          io.Discard,
		StackTraceLevel: riff.LevelError,
	})
	ctx := context.Background()

	for range 10_000_000 {
		log.Info(ctx, "Starting task",
			log.Any("device_unique_id", "G4000E-1000-F"),
			log.Any("task_id", 123456),
			log.Any("status", "success"),
			log.Any("template_name", "index.tpl"),
		)
	}
}
