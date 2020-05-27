package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type adjustSeg struct {
	b      float64
	e      float64
	factor float64
}

type adjustParam struct {
	vertical   []adjustSeg
	horizontal []adjustSeg
}

var adParam = adjustParam{
	vertical: []adjustSeg{
		{0, 0.4, 0.987},
		{0.4, 0.6, 1.005},
		{0.6, 1, 1.018},
	},
	horizontal: []adjustSeg{},
}

type config struct {
	test    bool
	srcFile string
	dstFile string
	srcDir  string
	dstDir  string
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

func test(cfg config) {
	adjust(cfg.srcFile, cfg.dstFile)
}

func parse() config {
	var cfg config
	flag.BoolVar(&cfg.test, "t", false, "test mode")
	flag.StringVar(&cfg.srcFile, "i", "", "source png")
	flag.StringVar(&cfg.dstFile, "o", "", "dest png")
	flag.StringVar(&cfg.srcDir, "sd", "", "source dir")
	flag.StringVar(&cfg.dstDir, "dd", "", "dest dir")
	flag.Parse()
	return cfg
}

func loop(cfg config) {
	delay := 1
	for {
		fn := waitForFile(cfg.srcDir)
		srcFile := filepath.Join(cfg.srcDir, fn)
		log.Printf("adjust %s file ...", srcFile)
		err := adjust(srcFile, filepath.Join(cfg.dstDir, fn))
		if err == nil {
			err = os.Remove(srcFile)
			delay = 1
		} else {
			log.Printf("adjust %s error: %v, retry after %d seconds...", srcFile, err, delay)
			time.Sleep(time.Duration(delay) * time.Second)
			if delay < 30 {
				delay *= 2
			}
		}
	}
}

func waitForFile(dir string) string {
	for {
		files, err := filepath.Glob(filepath.Join(dir, "*.png"))
		if err != nil {
			log.Printf("search png file error: %v", err)
			return ""
		}

		// Wait for file size not increase
		var curfs int64
	outer:
		for _, fn := range files {
			curfs = 0
			for {
				fs, err := os.Stat(fn)
				if err != nil {
					log.Printf("os.Stat(%s) error: %v", fn, err)
					continue outer
				}
				if curfs == fs.Size() {
					return filepath.Base(fn)
				}
				curfs = fs.Size()
				time.Sleep(time.Second)
			}
		}
		time.Sleep(1 * time.Second)
	}
}
