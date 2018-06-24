package epdownloader

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"github.com/meros/go-svtdownloader/libs/eplister"
	"github.com/meros/go-svtdownloader/libs/gqcommon"
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
func Get(episode eplister.Episode) error {
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

	folder := path.Join("videos", episode.Series, episode.Season)
	file := episode.Episode + ".mp4"
	fullPath := path.Join(folder, file)
	os.MkdirAll(folder, 0700)
	exec.Command("ffmpeg",
		"-i", videoURL,
		"-c", "copy",
		"-bsf:a", "aac_adtstoasc",
		fullPath).Run()

	return nil
}