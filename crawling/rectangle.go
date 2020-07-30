package main

import (
	"fmt"
	"strconv"
	"strings"
)

type rectangle struct {
	l float64
	b float64
	r float64
	t float64
}

func (r *rectangle) String() string {
	return fmt.Sprintf("%.4f,%.4f,%.4f,%.4f", r.l, r.b, r.r, r.t)
}

func (r *rectangle) Set(s string) error {
	var arrFloat []float64
	for _, vs := range strings.Split(s, ",") {
		f, err := strconv.ParseFloat(vs, 64)
		if err != nil {
			return err
		}
		arrFloat = append(arrFloat, f)
	}
	if len(arrFloat) != 4 {
		return fmt.Errorf("rectangle must be 4 float")
	}
	r.l, r.b, r.r, r.t = arrFloat[0], arrFloat[1], arrFloat[2], arrFloat[3]
	return nil
}
