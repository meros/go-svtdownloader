package main

import (
	"log"
	"os"
	"path"
	"regexp"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/meros/go-svtdownloader/libs/confreader"
	"github.com/meros/go-svtdownloader/libs/epdownloader"
	"github.com/meros/go-svtdownloader/libs/eplister"
	"github.com/meros/go-svtdownloader/libs/epnamer"
	pushbullet "github.com/xconstruct/go-pushbullet"
)

func main() {
	config := kingpin.Flag("config", "Config file").Short('c').Required().String()
	forever := kingpin.Flag("forever", "Keep running forever").Short('f').Bool()

	kingpin.Parse()

	file, err := os.Open(*config)
	if err != nil {
		log.Fatal("Failed to open config file: ", err)
	}

	mainConfig, err := confreader.Parse(file)
	if err != nil {
		log.Fatal("Failed to parse config file: ", err)
	}

	var pb *pushbullet.Client
	if mainConfig.PushbulletToken != "" {
		pb = pushbullet.New(mainConfig.PushbulletToken)
	}

	for {
		for _, serie := range mainConfig.Series {
			log.Println("Fetching series", serie.Key)
			eps, err := eplister.Get(serie.Key)
			if err != nil {
				log.Fatal("Failed to fetch serie", err)
			}

			epnamerOptions := epnamer.Options{
				Series: &epnamer.Replacement{
					Re:          *regexp.MustCompile(serie.Series.Regex),
					Replacement: serie.Series.Replacement},
				Season: &epnamer.Replacement{
					Re:          *regexp.MustCompile(serie.Season.Regex),
					Replacement: serie.Season.Replacement},
				Episode: &epnamer.Replacement{
					Re:          *regexp.MustCompile(serie.Episode.Regex),
					Replacement: serie.Episode.Replacement},
				TemplateString: serie.FilenameTemplate}

			for _, ep := range eps {
				filename, _ := epnamer.Filename(ep, epnamerOptions)
				filename = path.Join(mainConfig.BaseFolder, filename)

				err := epdownloader.Get(ep, filename)
				if err != nil {
					log.Fatal("Failed to download file", err)
				}

				if pb != nil {
					dev, err := pb.Device(mainConfig.PushbulletDevice)
					if err == nil {
						err = dev.PushNote("Episode downloaded", ep.Series+" "+ep.Season+" "+ep.Episode+" has been downloaded")
					}

					if err != nil {
						log.Println("Failed to push notification to device", err)
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
