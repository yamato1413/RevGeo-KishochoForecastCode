package common

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func ErrLog(err error) {
	if err != nil {
		log.Print(err)
	}
}

func GetJson(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte(""), err
	}
	res, err := http.DefaultClient.Do(req)
	ErrLog(err)

	body, err := ioutil.ReadAll(res.Body)
	ErrLog(err)
	defer res.Body.Close()

	return body, err
}

func Json2Map(body []byte) interface{} {
	var res interface{}
	json.Unmarshal(body, &res)
	return res
}
