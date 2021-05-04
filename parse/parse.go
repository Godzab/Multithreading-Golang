package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/minio/simdjson-go"
)

type Recipe struct {
	PostCode string `json:"postcode"`
	Recipe   string `json:"recipe"`
	Delivery string `json:"delivery"`
	Count    int    `json:"count"`
}

type BusiestPostcode struct {
	PostCode      string `json:"postcode"`
	DeliveryCount int    `json:"delivery_count"`
}

type CostPerPostalCodeAndTime struct {
	PostCode      string `json:"postcode"`
	From          string `json:"from"`
	To            string `json:"to"`
	DeliveryCount int    `json:"delivery_count"`
}

type Results struct {
	UniqueRecipeCount int `json:"unique_recipe_count"`
	CountPerRecipe    []struct {
		Recipe string `json:"recipe"`
		Count  int    `json:"count"`
	} `json:"count_per_recipe"`
	BusiestPostcode         BusiestPostcode          `json:"busiest_postcode"`
	CountPerPostcodeAndTime CostPerPostalCodeAndTime `json:"count_per_postcode_and_time"`
	MatchByName             []string                 `json:"match_by_name"`
}

var (
	lockRc = sync.Mutex{}
	lockPc = sync.Mutex{}
	lockSc = sync.Mutex{}

	timeLayout        = "3PM"
	UnqRecipeCnt      = 0
	searchResults     = []string{}
	queryX            = []string{"Potato", "Veggie", "Mushroom"}
	RecipeCount       = map[string]Recipe{}
	PostalCodeDC      = map[string]BusiestPostcode{}
	busiestPostalCode = BusiestPostcode{}
	PCAndTimeSch      = CostPerPostalCodeAndTime{}
	results           = Results{}
	fl                *os.File
	err               error
	pcQuery           string
	fromQuery         string
	toQuery           string

	timeRangeRegex = regexp.MustCompile(`([1-9]|1[0-2])([AP][M])\s-\s([1-9]|1[0-2])([AP][M])`)
	fromValidator  = regexp.MustCompile(`([1-9]|1[0-2])([AP][M])\s`)
	toValidator    = regexp.MustCompile(`([1-9]|1[0-2])([AP][M])$`)
)

func init() {
	source := flag.String("source", "", "Absolute path to source file.")
	search := flag.String("search", "Mushroom,Veggie,Potato", "Comma delimited recipe names to search.")
	pcAndTime := flag.String("postcodeTime", "10213:10AM:4PM", "Postcode with start time and end time colon demacated.")
	flag.Parse()

	if *source != "" {
		fl, err = os.Open(*source)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
	}

	if *search != "" {
		data := strings.Split(*search, ",")
		queryX = data
	}

	if *pcAndTime != "" {
		d := strings.Split(*pcAndTime, ":")
		if len(d) != 3 {
			fmt.Fprint(os.Stderr, "Search flag expected 3 fields. Please check help for format guide.")
			os.Exit(1)
		}
		pcQuery, fromQuery, toQuery = d[0], d[1], d[2]
		PCAndTimeSch = CostPerPostalCodeAndTime{pcQuery, fromQuery, toQuery, 0}
	}
}

func parseJson() {
	// Temp values.
	//r, _ := ioutil.ReadFile("/Users/godfreybafana/projects/golang/assessment/sample.json")
	iter, _ := simdjson.Parse(fl)
	for {
		typ := iter.Advance()

		switch typ {
		case simdjson.TypeRoot:
			if typ, tmp, err = iter.Root(tmp); err != nil {
				return
			}

			if typ == simdjson.TypeObject {
				if obj, err = tmp.Object(obj); err != nil {
					return
				}

				e := obj.FindKey(key, &elem)
				if e != nil && elem.Type == simdjson.TypeString {
					v, _ := elem.Iter.StringBytes()
					fmt.Println(string(v))
				}
			}

		default:
			return
		}
	}
}

func GenerateRecipeRpt() {
	dftStart, _ := time.Parse(timeLayout, strings.ReplaceAll(fromQuery, " ", ""))
	dftEnd, _ := time.Parse(timeLayout, strings.ReplaceAll(toQuery, " ", ""))
	decoder := json.NewDecoder(fl)

	_, err := decoder.Token()
	if err != nil {
		fmt.Fprintf(os.Stdout, "%s\n", err.Error())
		os.Exit(2)
	}

	for decoder.More() {
		var rcp Recipe
		err := decoder.Decode(&rcp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(2)
		}

		if rcp.PostCode == pcQuery {
			updateCustomQuery(rcp, dftStart, dftEnd)
		}

		rcpSearchKey := strings.ReplaceAll(strings.ToLower(strings.TrimSpace(rcp.Recipe)), " ", "")
		pcSearchKey := strings.ReplaceAll(strings.ToLower(strings.TrimSpace(rcp.PostCode)), " ", "")
		go recipeCountOperation(rcp, rcpSearchKey, pcSearchKey)
		go searchStrings(rcp, queryX)
	}

	_, err = decoder.Token()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(2)
	}

	//Insert Mutex here
	for _, bpc := range PostalCodeDC {
		go func(bpc BusiestPostcode) {
			if busiestPostalCode.PostCode == "" {
				busiestPostalCode = bpc
			} else {
				if busiestPostalCode.DeliveryCount < bpc.DeliveryCount {
					busiestPostalCode = bpc
				}
			}
		}(bpc)
	}

	for _, rcp := range RecipeCount {
		results.CountPerRecipe = append(results.CountPerRecipe, struct {
			Recipe string `json:"recipe"`
			Count  int    `json:"count"`
		}{rcp.Recipe, rcp.Count})
	}

	results.UniqueRecipeCount = UnqRecipeCnt
	results.BusiestPostcode = busiestPostalCode
	results.CountPerPostcodeAndTime = PCAndTimeSch
	results.MatchByName = searchResults

	js, _ := json.Marshal(results)
	fmt.Fprintf(os.Stdout, "%s\n", js)
}

func recipeCountOperation(r Recipe, rcpSearchKey, pcSearchKey string) {
	// Count recipes
	lockRc.Lock()
	if rcp, rcpFnd := RecipeCount[rcpSearchKey]; rcpFnd {
		rcp.Count += 1
		RecipeCount[rcpSearchKey] = rcp
	} else {
		UnqRecipeCnt += 1
		r.Count = 1
		RecipeCount[rcpSearchKey] = r
	}
	lockRc.Unlock()

	// Rank deliveries by postal code
	lockPc.Lock()
	if pc, pcFnd := PostalCodeDC[pcSearchKey]; pcFnd {
		pc.DeliveryCount += 1
		PostalCodeDC[pcSearchKey] = pc

	} else {
		pc.DeliveryCount = 1
		pc.PostCode = r.PostCode
		PostalCodeDC[pcSearchKey] = pc
	}
	lockPc.Unlock()
}

func updateCustomQuery(r Recipe, dftStart, dftEnd time.Time) {
	if timeRangeRegex.MatchString(r.Delivery) {
		from := string(fromValidator.Find([]byte(r.Delivery)))
		to := string(toValidator.Find([]byte(r.Delivery)))
		start, _ := time.Parse(timeLayout, strings.ReplaceAll(from, " ", ""))
		end, _ := time.Parse(timeLayout, strings.ReplaceAll(to, " ", ""))

		if checkTimeRange(start, end, dftStart, dftEnd) {
			PCAndTimeSch.DeliveryCount += 1
		}
	}
}

func searchStrings(r Recipe, s []string) {
	//Implement Mutex here on routine
	for _, str := range s {
		if found := strings.Contains(strings.ToLower(r.Recipe), strings.ToLower(str)); found {
			lockSc.Lock()
			searchResults = append(searchResults, r.Recipe)
			sort.Strings(searchResults)
			lockSc.Unlock()
			return
		}
	}
}

func checkTimeRange(start, end, dftStart, dftEnd time.Time) bool {
	return (start.After(dftStart) || start.Equal(dftStart) &&
		(end.Before(dftEnd) || end.Equal(dftEnd)) &&
		end.After(start))
}
