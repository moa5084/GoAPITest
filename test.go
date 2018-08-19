package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	URL   = "https://map.yahooapis.jp/weather/V1/place"
	APPID = "dj00aiZpPVh2SFpIcXlCcUw0ayZzPWNvbnN1bWVyc2VjcmV0Jng9OGI-"
)

// リクエストパラメータを表現する構造体
type Param struct {
	Prop, Data string
}

// レスポンスXMLを表現する構造体色々
type ResultXML struct {
	Features []Feature `xml:"Feature"`
}

type Feature struct {
	//Name    string  `xml:"Name"`
	Geometry Geometry `xml:"Geometry"`
	Property Property `xml:"Property"`
}

type Geometry struct {
	Coordinates string `xml:"Coordinates"`
}

type Property struct {
	WeatherList WeatherList `xml:"WeatherList"`
}

type WeatherList struct {
	Weathers []Weather `xml:"Weather"`
}

type Weather struct {
	Type     string `xml:"Type"`
	Date     string `xml:"Date"`
	Rainfall string `xml:"Rainfall"`
}

// レスポンスJSONを表現する構造体色々
type Result struct {
	Name        string `json:"Feature.Name"`
	CountryName string `json:"Feature.Property.Country.Name"`
	Address     string `json:"Feature.Property.Address"`
}

func main() {
	// リクエストパラメータの設定
	q_appid := Param{"appid", APPID}
	q_query := Param{"coordinates", "35.680029,139.737236"}
	q_output := Param{"output", "xml"}
	//q_output := Param{"output", "json"}
	p := []Param{q_appid, q_query, q_output}
	// http Request を作成
	req, errReq := http.NewRequest("GET", URL+"?"+createQuery(p), nil)

	fmt.Println(URL + "?" + createQuery(p))

	if errReq != nil {
		// todo
	}

	// http Request
	client := new(http.Client)
	resp, errResp := client.Do(req)
	if errResp != nil {
		// todo
	}
	defer resp.Body.Close()

	// 結果を出力（JSON）
	byteArray, errRead := ioutil.ReadAll(resp.Body)
	if errRead == nil {
		// とりあえずフォーマットせずに出力
		/*
			var out bytes.Buffer
			json.Indent(&out, byteArray, "", "\t")
			out.WriteTo(os.Stdout)
		*/
		fmt.Println(string(byteArray))

		// XMLのデコード
		resultXML := new(ResultXML)
		if err := xml.Unmarshal(byteArray, resultXML); err != nil {
			// todo
		}
		// fmt.Printf("Name:\t%s\nCountry:\t%s\nAddress:\t%s", resultXML.Name, resultXML.CountryName, resultXML.Address)
		fmt.Printf("Geometry:\t%s\nType:\t%s\nDate:\t%s\nRainfall:\t%s mm/h\n", resultXML.Features[0].Geometry.Coordinates, resultXML.Features[0].Property.WeatherList.Weathers[0].Type, resultXML.Features[0].Property.WeatherList.Weathers[0].Date, resultXML.Features[0].Property.WeatherList.Weathers[0].Rainfall)
		// JSON のデコード（文字コードが違う？）
		/*
			var results Result
			if err := json.Unmarshal(byteArray, results); err != nil {
				// todo
			}
			fmt.Println(results)
			fmt.Printf("Name:\t%s\nCountry:\t%s\nAddress:\t%s", results.Name, results.CountryName, results.Address)
		*/
	}
}

func createQuery(p []Param) string {
	var s string = ""
	for i := range p {
		if s != "" {
			s = s + "&"
		}
		s = s + fmt.Sprintf("%s=%s", p[i].Prop, p[i].Data)
	}
	return s
}
