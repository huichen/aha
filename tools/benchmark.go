package main

import (
	core "../core"
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"time"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "处理器profile文件")
	memprofile = flag.String("memprofile", "", "内存profile文件")
	service    core.LookupService
)

func main() {
	flag.Parse()

	// 捕获ctrl-c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			log.Print("捕获Ctrl-c，退出服务器")
			os.Exit(0)
		}
	}()

	service.Init()

	// 写入内存profile文件
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		defer f.Close()
	}

	// 打开处理器profile文件
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	numQueries := 100 // 1M
	t1 := time.Now()
	for i := 0; i < numQueries; i++ {
		service.GetMetricStats([]string{"112987:918142"}, []uint32{22212}, true)
	}
	t2 := time.Now()
	t := t2.Sub(t1).Seconds()
	log.Printf("索引%d次，耗时%f秒，QPS=%f", numQueries, float64(t), float64(numQueries)/float64(t))
}
