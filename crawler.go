package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getNumberFromText(text string) (number int) {
	text = strings.Replace(text, ",", "", -1)
	number, _ = strconv.Atoi(text)
	return
}

type informations struct {
	Infected string
	Restore  string
	Die      string
}

func getTextFromHTML(body io.Reader) (bodyBytes []byte) {
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return nil
	}
	return
}

// GetNumbers returns the numbers of the patients
// with following keys:
// 'confirmed': numbers of people who got confirmed
// 'death': numbers of people who died because of the corona19
// 'cured' : numbers of people who cured from the corona19
func GetNumbers() map[string]int {
	const URL = "http://happycastle.club/status?country=%EB%8C%80%ED%95%9C%EB%AF%BC%EA%B5%AD"
	var info informations
	res, err := http.Get(URL)
	if err != nil {
		return nil
	}
	defer res.Body.Close()

	json.Unmarshal(getTextFromHTML(res.Body), &info)

	confirmed, _ := strconv.Atoi(info.Infected)
	death, _ := strconv.Atoi(info.Die)
	cured, _ := strconv.Atoi(info.Restore)
	var result = map[string]int{
		"confirmed": confirmed,
		"death":     death,
		"cured":     cured,
	}
	return result
}

// GetNumbersFromNaver returns the numbers of the patients
// with following keys:
// 'confirmed': numbers of people who got confirmed
// 'death': numbers of people who died because of the corona19
// 'cured' : numbers of people who cured from the corona19
func GetNumbersFromNaver() map[string]int {
	const URL = "https://m.search.naver.com/search.naver?query=%EC%BD%94%EB%A1%9C%EB%82%9819"
	const confirmedPatients = "#_cs_common_production > div > div.status_info > ul > li.info_01 > p"
	const curedPatients = "#_cs_common_production > div > div.status_info > ul > li.info_03 > p"
	const diedPatients = "#_cs_common_production > div > div.status_info > ul > li.info_04 > p"
	numbers := make(map[string]int)

	res, err := http.Get(URL)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil
	}

	doc.Find(confirmedPatients).Each(func(i int, s *goquery.Selection) {
		numbers["confirmed"] = getNumberFromText(s.Text())
	})
	doc.Find(diedPatients).Each(func(i int, s *goquery.Selection) {
		numbers["death"] = getNumberFromText(s.Text())
	})
	doc.Find(curedPatients).Each(func(i int, s *goquery.Selection) {
		numbers["cured"] = getNumberFromText(s.Text())
	})

	return numbers
}
