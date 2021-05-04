package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	windRegex     = regexp.MustCompile(`\d* METAR.*EGLL \d*Z [A-Z]*(\d{5}KT|VRB\d{2}KT).*=`)
	tafValidation = regexp.MustCompile(`.*TAF.*`)
	comment       = regexp.MustCompile(`\w*#.*`)
	metarClose    = regexp.MustCompile(`.*=`)
	variableWind  = regexp.MustCompile(`.*VRB\d{2}KT`)
	validWind     = regexp.MustCompile(`\d{5}KT`)
	windDirOnly   = regexp.MustCompile(`(\d{3})\d{2}KT`)
	windDist      [8]int
)

func main() {
	textChannel := make(chan string)
	metarChannel := make(chan []string)
	windsChannel := make(chan []string)
	resultsChannel := make(chan [8]int)

	go parseToArray(textChannel, metarChannel)
	go extractWindDirection(metarChannel, windsChannel)
	go mineWindDistribution(windsChannel, resultsChannel)

	abspath, _ := filepath.Abs("metarfiles/")
	files, _ := ioutil.ReadDir(abspath)
	fmt.Println(len(files))
	start := time.Now()
	for _, file := range files {
		dat, err := ioutil.ReadFile(filepath.Join(abspath, file.Name()))
		if err != nil {
			panic(err)
		}
		text := string(dat)
		textChannel <- text
	}
	close(textChannel)
	results := <-resultsChannel
	elapsed := time.Since(start)
	fmt.Printf("%v\n", results)
	fmt.Printf("Time taken was %s\n", elapsed)
}

func parseToArray(txt chan string, metar chan []string) {
	for data := range txt {
		lines := strings.Split(data, "\n")
		metarSlice := make([]string, len(lines))
		metarStr := ""

		for _, ln := range lines {
			if tafValidation.MatchString(ln) {
				break
			}
			if !comment.MatchString(ln) {
				metarStr += strings.Trim(ln, " ")
			}
			if metarClose.MatchString(ln) {
				metarSlice = append(metarSlice, metarStr)
				metarStr = ""
			}
		}
		metar <- metarSlice
	}
	close(metar)
}

func extractWindDirection(metarChannel chan []string, windsChannel chan []string) {

	for data := range metarChannel {
		winds := make([]string, 0, len(data))
		for _, metar := range data {
			if windRegex.MatchString(metar) {
				winds = append(winds, windRegex.FindAllStringSubmatch(metar, -1)[0][1])
			}
		}
		windsChannel <- winds
	}
	close(windsChannel)
}

func mineWindDistribution(windsChannel chan []string, results chan [8]int) {
	for winds := range windsChannel {
		for _, wind := range winds {
			if variableWind.MatchString(wind) {
				for i := 0; i < 8; i++ {
					windDist[1]++
				}
			} else if validWind.MatchString(wind) {
				windStr := windDirOnly.FindAllStringSubmatch(wind, -1)[0][1]
				if d, err := strconv.ParseFloat(windStr, 64); err == nil {
					dirIndex := int(math.Round(d/45.0)) % 8
					windDist[dirIndex]++
				}
			}
		}
	}
	results <- windDist
	close(results)
}
