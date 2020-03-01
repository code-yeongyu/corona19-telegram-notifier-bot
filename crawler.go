package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getNumberFromText(text string) (number int) {
	text = text[:len(text)-4]
	text = strings.Replace(text, ",", "", -1)
	number, _ = strconv.Atoi(text)
	return
}

// GetNumbers returns the numbers of the patients
// with following keys:
// 'confirmed': numbers of people who got confirmed
// 'death': numbers of people who died because of the corona19
// 'cured' : numbers of people who cured from the corona19
func GetNumbers() map[string]int {
	const URL = "http://ncov.mohw.go.kr/index_main.jsp"
	const confirmedPatients = "body > div > div.container.main_container > div > div:nth-child(1) > div.co_cur > ul > li:nth-child(1) > a"
	const curedPatients = "body > div > div.container.main_container > div > div:nth-child(1) > div.co_cur > ul > li:nth-child(2) > a"
	const diedPatients = "body > div > div.container.main_container > div > div:nth-child(1) > div.co_cur > ul > li:nth-child(3) > a"
	numbers := make(map[string]int)

	res, err := http.Get(URL)
	if err != nil {
		return nil
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

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
