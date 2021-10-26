# 緯度経度から気象庁の地域コードを取得する

2021年4月より、気象庁の天気予報がJSONとして取得できるようになりました。
天気予報を取得したい地域の地域コードでリクエストを送ると、概要や3日間天気、週間天気などが取得できます。

しかしこの地域コード、総務省が管理している自治体コードがベースになっており、あまりなじみがありません。さらに、気象庁が公開している地域コード一覧JSONは非常に扱いにくい構造になっており、目的の地域コードを探すのも大変です。

そこで、緯度経度から
* 天気予報リクエスト用コード
* コードに対応する市町村名
* 各レベル別地域コード

を取得できるようにしてみました。


```
https://revgeo-forecastcode.herokuapp.com/lat={lat}+lon={lon}
```

上記のURLの{lat}を緯度に、{lon}を経度に置き換えてhttpリクエストを送ると、以下のようなJSONを返します。

（例）京都御苑の緯度経度でリクエスト
```
https://revgeo-forecastcode.herokuapp.com/lat=35.021077+lon=135.761731
```
```
{
  "forecastcode": "260000",
  "officename": "京都府",
  "cityname": "京都市",
  "centers": "010600",
  "offices": "260000",
  "class10s": "260010",
  "class15s": "260011",
  "class20s": "2610000"
}
```
ここで取得できたforecastcodeを利用すれば、気象庁から天気予報を取得できます。

（例）京都府の天気概要
```
https://www.jma.go.jp/bosai/forecast/data/overview_forecast/260000.json
```
（例）京都府の3日間天気
```
https://www.jma.go.jp/bosai/forecast/data/forecast/260000.json
```

3日間天気のデータにはその都道府県内の各地域の天気予報が含まれています。希望の地域のデータを抽出するためにclass10sのコード等を使用します。

なお、緯度経度から自治体コードへの変換は国土地理院の逆ジオコーディングAPIを利用しました。
```
https://mreversegeocoder.gsi.go.jp/reverse-geocoder/LonLatToAddress?lat={lat}&lon={lon}
```