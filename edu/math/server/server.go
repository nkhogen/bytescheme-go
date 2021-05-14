package main

import (
    "fmt"
    "net/http"
	"bytescheme/edu/math"
	"strconv"
)

func handler(w http.ResponseWriter, r *http.Request) {
	gen := &math.ProblemGenerator{Conf: &math.Config{Cols: 2}}
	gen.Conf.MultConf = &math.MultConfig{FirstLen: 3, SecondLen: 2, Size: 5}
	gen.Conf.DivConf = &math.DivConfig{QuotientLen: 3, DivisorLen: 1, Size: 5}
	gen.Conf.SubConf = &math.SubConfig{SubtractorLen: 3, ResultLen: 3, Size: 5}
	query := r.URL.Query()
	multSizeParam := query.Get("multsize")
	if multSizeParam != "" {
		multSize, err := strconv.Atoi(multSizeParam)
		if err == nil {
			gen.Conf.MultConf.Size = multSize
		}
	}
	divSizeParam := query.Get("divsize")
	if divSizeParam != "" {
		divSize, err := strconv.Atoi(divSizeParam)
		if err == nil {
			gen.Conf.DivConf.Size = divSize
		}
	}
	subSizeParam := query.Get("divsize")
	if subSizeParam != "" {
		subSize, err := strconv.Atoi(subSizeParam)
		if err == nil {
			gen.Conf.SubConf.Size = subSize
		}
	}
	html := gen.GenerateHTML()
	w.Header().Set("Content-Type", "text/html")
    fmt.Fprintf(w, "%s\n", string(html))
}

func main() {
    http.HandleFunc("/v1/problems", handler)
    http.ListenAndServe(":9090", nil)
}