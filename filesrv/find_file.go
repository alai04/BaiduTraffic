package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type stFile struct {
	Time     string `json:"time"`
	Filename string `json:"filename"`
}

const (
	fileDateTimeFormat = "20060102_150405"
	fileNameFormat     = "L%d_%s.png"
)

type fileFinder struct {
	level     int
	tm        time.Time
	deviation int
}

func (ff fileFinder) findFile() (res stFile, err error) {
	fn, tm := ff.findNearestFile()
	if fn == "" {
		return res, fmt.Errorf("cannot find image")
	}
	res.Time = tm.Format(dateTimeFormat)
	res.Filename = fn
	return
}

func (ff fileFinder) findNearestFile() (fn string, ftm time.Time) {
	diffMin := float64(time.Duration(1<<63 - 1))
	filePrefix := fmt.Sprintf("L%d_", ff.level)
	filepath.Walk(fileDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		fnCur := info.Name()
		if !info.IsDir() && strings.HasPrefix(fnCur, filePrefix) {
			tmCur, err := fileTime(fnCur)
			if err != nil {
				fmt.Printf("get time from filename %q error: %v\n", fnCur, err)
			}
			diff := math.Abs(float64(ff.tm.Sub(tmCur)))
			if diff < diffMin {
				diffMin = diff
				fn = fnCur
				ftm = tmCur
			}
		}
		return nil
	})
	if diffMin > float64(time.Duration(ff.deviation)*time.Second) {
		return "", ftm
	}
	return fn, ftm
}

func fileTime(fn string) (tm time.Time, err error) {
	b := len(fn) - 19
	return time.Parse(fileDateTimeFormat, fn[b:b+15])
}
