package main

import (
	"log"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/meros/go-svtdownloader/libs/epdownloader"
	"github.com/meros/go-svtdownloader/libs/eplister"
	pushbullet "github.com/xconstruct/go-pushbullet"
)

func main() {
	series := kingpin.Flag("series", "Name of series").Short('s').Required().Strings()
	outDir := kingpin.Flag("outDir", "Base directory to put files").Short('o').Required().String()
	pushbulletToken := kingpin.Flag("pushbulletToken", "Pushbullet token for notifications").Short('p').String()
	pushbulletDevice := kingpin.Flag("pushbulletDevice", "Pushbullet device for notifications").Short('d').String()
	forever := kingpin.Flag("forever", "Keep running forever").Short('f').Bool()

	kingpin.Parse()

	var pb *pushbullet.Client
	if *pushbulletToken != "" {
		pb = pushbullet.New(*pushbulletToken)
	}

	for {
		for _, serie := range *series {
			log.Println("Fetching series", serie)
			eps, _ := eplister.Get(serie)
			for _, ep := range eps {
				err := epdownloader.Get(ep, *outDir)
				if err != nil {
					log.Println(err)
					continue
				}

				if pb != nil {
					dev, err := pb.Device(*pushbulletDevice)
					if err == nil {
						err = dev.PushNote("Episode downloaded", ep.Series+" "+ep.Season+" "+ep.Series+" has been downloaded")
					}

					if err != nil {
						log.Println("Failed to push notification to device", *pushbulletDevice, err)
					}
				}
			}
		}

		if !(*forever) {
			return
		}

		log.Println("Sleeping for 10 minutes and checking again at that point")
		time.Sleep(10 * time.Minute)
	}
}
