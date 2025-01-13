package main

import (
	"encoding/csv"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {
	thin, fat, experiment := getData()
	plot := charts.NewLine()
	xAxis := make([]int, len(thin))
	for i := range xAxis {
		xAxis[i] = i
	}

	plot.SetGlobalOptions(
		charts.WithXAxisOpts(opts.XAxis{Name: "loop count"}),
		charts.WithYAxisOpts(opts.YAxis{Name: "Nanosecond"}),
	)

	plot.
		SetXAxis(xAxis).
		AddSeries("Thin", thin).
		AddSeries("Fat", fat).
		AddSeries("Experiment", experiment)

	f, err := os.Create("plot.html")
	if err != nil {
		panic(err)
	}
	err = plot.Render(f)
	if err != nil {
		panic(err)
	}
}

func getData() ([]opts.LineData, []opts.LineData, []opts.LineData) {
	var thin []opts.LineData
	var fat []opts.LineData
	var experiment []opts.LineData

	f, err := os.Open("data")
	if err != nil {
		panic(err)
	}

	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		series := record[0]
		val, err := strconv.ParseUint(record[1], 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		point := opts.LineData{Value: val, Name: "ns"}

		switch series {
		case "thin":
			thin = append(thin, point)
		case "fat":
			fat = append(fat, point)
		case "experiment":
			experiment = append(experiment, point)
		}
	}

	return thin, fat, experiment
}
