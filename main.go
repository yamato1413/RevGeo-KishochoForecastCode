package main

import (
	"RevGeo-KishochoForecastCode/area"
	"RevGeo-KishochoForecastCode/common"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/koron/go-dproxy"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/lat={lat}+lon={lon}", procRequest)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}

type Results struct {
	ForecastCode string `json:"forecastcode"`
	PrefName     string `json:"prefname"`
	Cityname     string `json:"cityname"`
	Centers      string `json:"centers"`
	Offices      string `json:"offices"`
	Class10s     string `json:"class10s"`
	Class15s     string `json:"class15s"`
	Class20s     string `json:"class20s"`
}

func procRequest(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	lat, err := strconv.ParseFloat(vars["lat"], 64)
	common.ErrLog(err)
	lon, err := strconv.ParseFloat(vars["lon"], 64)
	common.ErrLog(err)
	localcode := latlonToLocalCode(lat, lon)

	res := new(Results)
	if localcode != "" {
		citycode := toCityCode(localcode)
		res.Class20s = citycode + "00"
		res.Class15s = parentcode(res.Class20s)
		res.Class10s = parentcode(res.Class15s)
		res.Offices = parentcode(res.Class10s)
		res.Centers = parentcode(res.Offices)
		switch res.Offices {
		// 十勝地方
		case "014030":
			res.ForecastCode = "014100" // 釧路地方
		// 奄美地方
		case "460040":
			res.ForecastCode = "460100" // 鹿児島県
		default:
			res.ForecastCode = res.Offices
		}
		res.PrefName = toPrefName(res.Offices)
		res.Cityname = toCityName(res.Class20s)
	}
	json.NewEncoder(w).Encode(res)
}

type RevGeo struct {
	Results struct {
		MuniCd string `json:"muniCd"`
		Lv01Nm string `json:"lv01Nm"`
	} `json:"results"`
}

func latlonToLocalCode(lat, lon float64) string {
	URL := fmt.Sprintf("https://mreversegeocoder.gsi.go.jp/reverse-geocoder/LonLatToAddress?lat=%v&lon=%v", lat, lon)
	body, err := common.GetJson(URL)
	common.ErrLog(err)

	var rg RevGeo
	json.Unmarshal(body, &rg)

	return rg.Results.MuniCd
}

func parentcode(code string) string {
	category := []string{"class20s", "class15s", "class10s", "offices"}
	areainfo := dproxy.New(area.AreaInfoMap())
	for _, class := range category {
		pcode, err := areainfo.M(class).M(code).M("parent").String()
		if err == nil {
			return pcode
		}
	}
	return ""
}

func toCityName(citycode string) string {
	areainfo := dproxy.New(area.AreaInfoMap())
	name, _ := areainfo.M("class20s").M(citycode).M("name").String()
	return name
}
func toPrefName(officecode string) string {
	areainfo := dproxy.New(area.AreaInfoMap())
	name, _ := areainfo.M("offices").M(officecode).M("name").String()
	return name
}

// 区レベルの自治体コードを市レベルに直す
// 二分探索で近似値検索をしている
func toCityCode(code string) string {
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
