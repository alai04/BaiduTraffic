package main

import (
	"flag"
	"log"
)

func main() {
	level := flag.Int("l", 9, "level of map")
	workers := flag.Int("w", 1, "num of workers")
	flag.Parse()
	m := NewTrafficMap(*level, *workers)
	fn, err := m.GetMap()
	if err == nil {
		log.Printf("Traffic map save to: %v\n", fn)
	} else {
		log.Printf("GetTrafficMap() error: %v", err)
	}
}
