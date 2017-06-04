package main

import (
	"bufio"
	"fmt"
	"image/color"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
)

var linerex = regexp.MustCompile(`^Benchmark.+/(\d+)-\d+\s+\d+\s+(\d+) ns/op\s+(\d+) B/op`)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("file require")
	}

	fn := os.Args[1]

	name := strings.Split(fn, ".")[0]

	dur, mem := getPoints(fn)

	{
		p, err := plot.New()
		if err != nil {
			log.Fatal(err)
		}

		p.Title.Text = name + " duration"
		p.X.Label.Text = "10^6 Links"
		p.Y.Label.Text = "ms"

		durLine, err := plotter.NewLine(dur)
		if err != nil {
			log.Fatal(err)
		}
		durLine.Color = color.RGBA{B: 255, A: 255}

		p.Add(durLine)
		if err := p.Save(12*vg.Inch, 6*vg.Inch, name+" duration.png"); err != nil {
			panic(err)
		}
	}
	{
		p, err := plot.New()
		if err != nil {
			log.Fatal(err)
		}

		p.Title.Text = name + " memory allocated"
		p.X.Label.Text = "10^6 Links"
		p.Y.Label.Text = "MB"

		memLine, err := plotter.NewLine(mem)
		if err != nil {
			log.Fatal(err)
		}
		memLine.Color = color.RGBA{G: 255, A: 255}

		p.Add(memLine)
		if err := p.Save(12*vg.Inch, 6*vg.Inch, name+" memory.png"); err != nil {
			panic(err)
		}
	}
}

func getPoints(fn string) (dur, mem plotter.XYs) {
	f, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		sub := linerex.FindStringSubmatch(scanner.Text())
		if len(sub) == 4 {
			var m, t, b float64

			fmt.Sscan(sub[1], &m)
			fmt.Sscan(sub[2], &t)
			fmt.Sscan(sub[3], &b)

			dur = append(dur, struct{ X, Y float64 }{m / 1000000, t / 1000000})
			mem = append(mem, struct{ X, Y float64 }{m / 1000000, b / (1024 * 1024)})
		}
	}

	return
}
