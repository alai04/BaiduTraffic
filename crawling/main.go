package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

const (
	defaultWorkers  = 10
	defaultMaxLevel = 11
	tmpFilename     = ".__tmp__.png"
)

type config struct {
	test    bool
	level   int
	workers int
	fileDir string
}

func main() {
	cfg := parse()
	if cfg.test {
		test(cfg)
	} else {
		fmt.Println("Run as deamon, Crtl-C to exit ...")
		loop(cfg)
	}
}

func parse() config {
	var cfg config
	flag.BoolVar(&cfg.test, "t", false, "test mode")
	flag.IntVar(&cfg.level, "l", 9, "max level of map, or level of map for test")
	flag.IntVar(&cfg.workers, "w", 1, "num of workers")
	flag.StringVar(&cfg.fileDir, "d", "./", "directory of output")
	flag.Parse()
	return cfg
}

func test(cfg config) {
	m := NewTrafficMap(cfg.level, cfg.workers, cfg.fileDir)
	fn, err := m.GetMap()
	if err == nil {
		log.Printf("Traffic map save to: %v", fn)
	} else {
		log.Printf("GetTrafficMap() error: %v", err)
	}
}

func loop(cfg config) {
	workers := cfg.workers
	if workers < 1 || workers > 20 {
		workers = defaultWorkers
	}
	maxLevel := cfg.level
	if maxLevel < 9 || workers > 15 {
		maxLevel = defaultMaxLevel
	}

	for {
		tBegin := time.Now()
		tNext := tBegin.Add(time.Minute)
		for l := 9; l <= maxLevel; l++ {
			m := NewTrafficMap(l, workers, cfg.fileDir)
			fn, err := m.GetMap()
			if err == nil {
				log.Printf("Traffic map save to: %v", fn)
			} else {
				log.Printf("GetTrafficMap() error: %v", err)
			}
		}
		log.Printf("Loop consuming: %v", time.Since(tBegin))
		if time.Now().Before(tNext) {
			d := time.Until(tNext)
			log.Printf("Wait %v for next loop", d)
			time.Sleep(d)
		}
	}
}
