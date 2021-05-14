package main

import (
	"bytescheme/edu/math"
)

func main() {
	gen := &math.ProblemGenerator{Conf: &math.Config{}}
	gen.Conf.MultConf = &math.MultConfig{FirstLen: 3, SecondLen: 2, Size: 5}
	gen.Conf.DivConf = &math.DivConfig{QuotientLen: 3, DivisorLen: 1, Size: 5}
	gen.Conf.SubConf = &math.SubConfig{SubtractorLen: 3, ResultLen: 3, Size: 5}
	err := gen.GenerateHTMLFile("/Users/naorem.khogendro.singh/Documents/problems-maths.html")
	if err != nil {
		panic(err)
	}
}