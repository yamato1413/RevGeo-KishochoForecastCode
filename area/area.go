package area

import (
	"RevGeo-KishochoForecastCode/common"
)

var areainfomap interface{}

func AreaInfoMap() interface{} {
	if areainfomap != nil {
		return areainfomap
	}
	body, err := common.GetJson("https://www.jma.go.jp/bosai/common/const/area.json")
	common.ErrLog(err)
	areainfomap = common.Json2Map(body)
	return areainfomap
}
