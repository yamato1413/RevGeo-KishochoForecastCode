package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"rev-kishocho-code/area"
	"rev-kishocho-code/common"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/koron/go-dproxy"
)

type Geo struct {
	Results struct {
		MuniCd string `json:"muniCd"`
		Lv01Nm string `json:"lv01Nm"`
	} `json:"results"`
}

func latlon_to_localcode(lat, lon float64) string {
	url := fmt.Sprintf("https://mreversegeocoder.gsi.go.jp/reverse-geocoder/LonLatToAddress?lat=%v&lon=%v", lat, lon)
	body, err := common.GetJson(url)
	common.ErrLog(err)

	var g Geo
	json.Unmarshal(body, &g)

	return g.Results.MuniCd
}

func parentcode(code string) (string, error) {
	category := []string{"class20s", "class15s", "class10s", "offices"}
	areainfo := dproxy.New(area.AreaInfoMap())
	for _, class := range category {
		pcode, err := areainfo.M(class).M(code).M("parent").String()
		if err == nil {
			return pcode, nil
		}
	}
	return "", errors.New("NotFound")
}

func localcode_to_citycode(code string) string {
	// 区レベルの自治体コードを市レベルに直す
	// 二分探索で近似値検索をしている
	codes := loadCityCodes()
	left, right := 0, len(codes)-1
	var mid int
	for right-left > 1 {
		mid = (left + right) / 2
		if code == codes[mid] {
			break
		} else if code > codes[mid] {
			left = mid
		} else {
			right = mid
		}
	}
	if code < codes[mid] {
		mid--
	}
	return codes[mid]
}

var cityCodes []string

// 総務省が公開している市レベル自治体コード表をCSVにしたものを読み込む
func loadCityCodes() []string {
	if len(cityCodes) != 0 {
		return cityCodes
	}
	file, err := os.Open("./code.csv")
	common.ErrLog(err)
	defer file.Close()

	reader := csv.NewReader(file)
	for {
		line, err := reader.Read()
		if err != nil {
			break
		}
		cityCodes = append(cityCodes, line[0])
	}
	return cityCodes
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/lat:{lat}lon:{lon}", procRequest)
	log.Fatal(http.ListenAndServe(":8081", r))
}

func procRequest(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	lat, _ := strconv.ParseFloat(vars["lat"], 64)
	lon, _ := strconv.ParseFloat(vars["lon"], 64)
	localcode := latlon_to_localcode(lat, lon)
	citycode := localcode_to_citycode(localcode)

	var res results
	res.Class20s = citycode + "00"
	res.Class15s, _ = parentcode(res.Class20s)
	res.Class10s, _ = parentcode(res.Class15s)
	res.Offices, _ = parentcode(res.Class10s)
	res.Centers, _ = parentcode(res.Offices)
	switch res.Offices {
	case "014030":
		res.ForecastRequest = "014100"
	case "460040":
		res.ForecastRequest = "460100"
	default:
		res.ForecastRequest = res.Offices
	}
	json.NewEncoder(w).Encode(res)
}

type results struct {
	Centers         string `json:"centers"`
	Offices         string `json:"offices"`
	Class10s        string `json:"class10s"`
	Class15s        string `json:"class15s"`
	Class20s        string `json:"class20s"`
	ForecastRequest string `json:"forecastrequest"`
}
