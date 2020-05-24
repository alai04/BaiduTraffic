package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	envWorkers  = "CRAWLING_WORKERS"
	envMaxLevel = "CRAWLING_MAXLEVEL"
)

type config struct {
	test    bool
	level   int
	workers int
}

func main() {
	cfg := parse()
	if cfg.test {
		test(cfg)
	} else {
		fmt.Println("Run as deamon, Crtl-C to exit ...")
		loop()
	}
}

func parse() config {
	var cfg config
	flag.BoolVar(&cfg.test, "t", false, "test mode")
	flag.IntVar(&cfg.level, "l", 9, "level of map")
	flag.IntVar(&cfg.workers, "w", 1, "num of workers")
	flag.Parse()
	return cfg
}

func test(cfg config) {
	m := NewTrafficMap(cfg.level, cfg.workers)
	fn, err := m.GetMap()
	if err == nil {
		log.Printf("Traffic map save to: %v", fn)
	} else {
		log.Printf("GetTrafficMap() error: %v", err)
	}
}

func loop() {
	workers, _ := strconv.Atoi(os.Getenv(envWorkers))
	if workers < 1 || workers > 20 {
		workers = 10
	}
	maxLevel, _ := strconv.Atoi(os.Getenv(envMaxLevel))
	if maxLevel < 9 || workers > 15 {
		maxLevel = 11
	}

	for {
		tBegin := time.Now()
		tNext := tBegin.Add(time.Minute)
		for l := 9; l <= maxLevel; l++ {
			m := NewTrafficMap(l, workers)
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
