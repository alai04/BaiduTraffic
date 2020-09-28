package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	defaultWorkers  = 10
	defaultMaxLevel = 11
	tmpFilename     = ".__tmp__.png"
)

type config struct {
	test     bool
	amap     bool
	level    int
	workers  int
	interval int
	fileDir  string
	keyFile  string
	rect     rectangle
}

var getRandomKey func() string

func main() {
	cfg := parse()

	keys, err := readKeys(cfg.keyFile)
	if err != nil {
		log.Fatal(err)
	}
	getRandomKey = func() string {
		return keys[rand.Intn(len(keys))]
	}

	if cfg.test {
		test(cfg)
	} else {
		fmt.Println("Run as deamon, Crtl-C to exit ...")
		loop(cfg)
	}
}

func parse() config {
	var cfg config
	flag.BoolVar(&cfg.test, "t", false, "测试模式，仅执行一次")
	flag.BoolVar(&cfg.amap, "a", false, "使用高德地图，缺省使用百度地图")
	flag.IntVar(&cfg.level, "l", 9, "抓取瓦片图的级别")
	flag.IntVar(&cfg.workers, "w", 1, "并发线程数")
	flag.IntVar(&cfg.interval, "i", 1, "间隔分钟数")
	flag.StringVar(&cfg.fileDir, "d", "./", "图片输出目录")
	flag.StringVar(&cfg.keyFile, "k", "keys.txt", "高德keys文件")
	flag.Var(&cfg.rect, "r", `抓取范围的坐标，格式：左下右上，如"121.4595,31.1939,121.5150,31.2503"`)
	flag.Parse()
	return cfg
}

func test(cfg config) {
	m := NewTrafficMap(cfg)
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

	for {
		tBegin := time.Now()
		tNext := tBegin.Add(time.Duration(rand.Int63n(int64(time.Duration(cfg.interval) * time.Minute))))
		m := NewTrafficMap(cfg)
		fn, err := m.GetMap()
		if err == nil {
			log.Printf("Traffic map save to: %v", fn)
		} else {
			log.Printf("GetTrafficMap() error: %v", err)
		}
		log.Printf("Loop consuming: %v", time.Since(tBegin))
		if time.Now().Before(tNext) {
			d := time.Until(tNext)
			log.Printf("Wait %v for next loop", d)
			time.Sleep(d)
		}
	}
}

func readKeys(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		key := scanner.Text()
		if len(key) == 32 {
			lines = append(lines, key)
		}
	}
	if len(lines) == 0 {
		return lines, fmt.Errorf("no valid key in %s", path)
	}
	return lines, scanner.Err()
}
