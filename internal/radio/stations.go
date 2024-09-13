package radio

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const defaultStationsCSV = `
	Chillsky,https://lfhh.radioca.st/stream,Chillhop;Downtempo
    Rainwave All,https://rainwave.cc/tune_in/5.mp3.m3u,Game Tracks and Remixes
    Rainwave Game,https://rainwave.cc/tune_in/1.mp3.m3u,Game Tracks
    Rainwave Chiptunes,https://rainwave.cc/tune_in/4.mp3.m3u,Chiptunes
    Rainwave OC Remix,https://rainwave.cc/tune_in/2.mp3.m3u,OverClocked Remixes
    Rainwave Covers,https://rainwave.cc/tune_in/3.mp3.m3u,Official and Fan Covers
	HiRes: Radio Paradise,https://stream.radioparadise.com/flacm,Eclectic
	HiRes: Radio Paradise Mellow,https://stream.radioparadise.com/mellow-flacm,Ballads;Eclectic
	HiRes: Radio Paradise Rock,https://stream.radioparadise.com/rock-flacm,Rock
	HiRes: SuperStereo 6: Instrumental Music,http://icecast.centaury.cl:7570/SuperStereoHiRes6,Instrumental
`

type Station struct {
	title string
	url   string
	tags  []string
}

func Stations(sta string) []Station {
	var reader *csv.Reader

	if sta == "" {
		reader = csv.NewReader(strings.NewReader(defaultStationsCSV))
	} else if strings.HasPrefix(sta, "http") {
		s, err := fetchStations(sta)
		if err != nil {
			s = defaultStationsCSV
		}
		reader = csv.NewReader(strings.NewReader(s))
	} else {
		file, err := os.Open(sta)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()
		reader = csv.NewReader(file)
	}

	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()

	if err != nil {
		fmt.Println("Can't parse stations CSV:", err)
		os.Exit(1)
	}

	stations := make([]Station, 0)
	for _, r := range records {
		s := new(Station)
		s.title = strings.Trim(r[0], " 	")
		s.url = strings.Trim(r[1], " 	")
		if len(r) > 2 {
			s.tags = strings.Split(strings.Trim(r[2], " 	"), ";")
		}
		stations = append(stations, *s)
	}

	return stations
}

func fetchStations(url string) (string, error) {
	resp, err := http.Get(url)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Failed to GET %s, status code: %d", url, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
