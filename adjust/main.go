package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/prometheus/client_golang/prometheus"
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
	port    string
}

func main() {
	cfg := parse()
	if cfg.test {
		test(cfg)
	} else {
		go metricsServe(cfg.port)
		go statFileNum(cfg.dstDir)

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
	flag.StringVar(&cfg.srcDir, "sd", "srcmap", "source dir")
	flag.StringVar(&cfg.dstDir, "dd", "dstmap", "dest dir")
	flag.StringVar(&cfg.port, "p", "8080", "port for metrics")
	flag.Parse()
	return cfg
}

func loop(cfg config) {
	delay := 1
	for {
		fn := waitForFile(cfg.srcDir)
		level := fn[1 : len(fn)-20]
		srcFile := filepath.Join(cfg.srcDir, fn)
		dstFile := filepath.Join(cfg.dstDir, fn)
		log.Printf("adjust %s file ...", srcFile)
		stage, err := adjust(srcFile, dstFile)
		if err == nil {
			err = os.Remove(srcFile)
			delay = 1
			adjustMapsTotal.With(prometheus.Labels{"level": level}).Inc()
			adjustConsumingTotalSeconds.With(prometheus.Labels{"level": level}).Observe(
				float64(timeConsuming) / float64(time.Second))
			fInfo, _ := os.Stat(dstFile)
			adjustTotalBytes.Observe(float64(fInfo.Size()))
		} else {
			adjustMapsFailure.With(prometheus.Labels{"level": level, "stage": stage}).Inc()
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
		files, err := filepath.Glob(filepath.Join(dir, "L*.png"))
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

func statFileNum(dir string) {
	const fileNameFormat = "L*_*.png"
	dir1 := fmt.Sprintf("%s/%s", dir, fileNameFormat)
	dir2 := fmt.Sprintf("%s/*/%s", dir, fileNameFormat)

	for {
		matches1, err1 := filepath.Glob(dir1)
		matches2, err2 := filepath.Glob(dir2)
		total := len(matches1) + len(matches2)
		if err1 == nil || err2 == nil {
			totalFilesInDest.Set(float64(total))
			log.Printf("Map files in %s: %d", dir, total)
		} else {
			err := err1
			if err == nil {
				err = err2
			}
			log.Printf("Glob(%q) return error: %v", dir, err)
		}
		time.Sleep(60 * time.Second)
	}
}
