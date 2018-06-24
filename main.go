package main

import (
	"log"

	"github.com/alecthomas/kingpin"
	"github.com/meros/go-svtdownloader/libs/epdownloader"
	"github.com/meros/go-svtdownloader/libs/eplister"
)

func main() {
	series := kingpin.Flag("series", "Name of series").Short('s').Required().String()
	outDir := kingpin.Flag("ourDir", "Base directory to put files").Short('o').Required().String()
	kingpin.Parse()

	eps, _ := eplister.Get(*series)
	for _, ep := range eps {
		err := epdownloader.Get(ep, *outDir)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
