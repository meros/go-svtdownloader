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
func Get(episode eplister.Episode, outDir string) error {
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

	folder := path.Join(outDir, episode.Series)
	fileBase := episode.Series + " " + episode.Season + " " + episode.Episode

	// IF file is on format Säsong X Avsnitt Y then change to SXEY for easier parsing
	re = regexp.MustCompile(`(^.*)Säsong ([0-9]+) Avsnitt ([0-9]+)(.*$)`)
	fileBase = re.ReplaceAllString(fileBase, `${1} S${2}E${3}${4}`)

	fileTemp := fileBase + ".part.mp4"
	file := fileBase + ".mp4"
	fullPath := path.Join(folder, file)
	fullPathTemp := path.Join(folder, fileTemp)

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
	os.MkdirAll(folder, 0777)
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
