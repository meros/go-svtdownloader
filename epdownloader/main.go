package epdownloader

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/meros/go-svtdownloader/eplister"
	"github.com/meros/go-svtdownloader/gqcommon"
)

const urlVideoApi = "http://api.svt.se/videoplayer-api/video/"

type jsonDataContainer struct {
	VideoPage struct {
		Video struct {
			ProgramTitle     string
			ProgramVersionID string
			ID               string
			Versions         []struct {
				ContentURL string
			}
		}
	}
}

type videoInfoContainer struct {
	VideoReferences []struct {
		URL    string
		Format string
	}
}

// Get will try to download video of provided episode
func Get(episode eplister.Episode, filename string) error {
	doc, err := gqcommon.GetDoc(episode.Url)
	if err != nil {
		return err
	}

	html, err := doc.Html()
	if err != nil {
		return err
	}

	re, err := regexp.Compile("__svtplay'] = ({.*});")
	if err != nil {
		return err
	}

	match := re.FindStringSubmatch(html)
	if len(match) != 2 {
		return errors.New("Wrong number of matches in trying to find json data")
	}

	jsonString := match[1]
	jsonData := &jsonDataContainer{}
	err = json.Unmarshal([]byte(jsonString), jsonData)
	if err != nil {
		log.Println(err)
		return err
	}

	videoInfoURL := urlVideoApi + jsonData.VideoPage.Video.ID
	videoInfoResponse, err := http.Get(videoInfoURL)
	if err != nil {
		log.Println(err)
		return err
	}

	videoInfoBytes, err := ioutil.ReadAll(videoInfoResponse.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	videoInfo := &videoInfoContainer{}
	err = json.Unmarshal(videoInfoBytes, videoInfo)
	if err != nil {
		log.Println(err)
		return err
	}

	var videoURL string
	for _, videoReference := range videoInfo.VideoReferences {
		if strings.Compare(videoReference.Format, "hls") != 0 {
			continue
		}

		videoURL = videoReference.URL
	}

	if videoURL == "" {
		return errors.New("Could not find suitable path to download")
	}

	// Fix broken thunderbirds series title
	fullPathTemp := filename + ".part.mp4"
	fullPath := filename + ".mp4"

	_, err = os.Stat(fullPath)
	if !os.IsNotExist(err) {
		log.Println("Not redownloading", fullPath)
		return errors.New("Episode already downloaded")
	}

	_, err = os.Stat(fullPathTemp)
	if !os.IsNotExist(err) {
		log.Println("Deleting old temp", fullPathTemp)
		os.Remove(fullPathTemp)
	}

	log.Println("Downloading", fullPath)
	os.MkdirAll(filepath.Dir(fullPath), 0777)
	err = exec.Command("ffmpeg",
		"-i", videoURL,
		"-c", "copy",
		"-bsf:a", "aac_adtstoasc",
		fullPathTemp).Run()

	if err != nil {
		return err
	}

	return os.Rename(fullPathTemp, fullPath)
}
