package main

import (
	"log"

	"github.com/meros/go-svtdownloader/libs/epdownloader"
	"github.com/meros/go-svtdownloader/libs/eplister"
)

func main() {
	eps, _ := eplister.Get("thunderbirds")
	for _, ep := range eps {
		err := epdownloader.Get(ep)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
