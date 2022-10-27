package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type kanji struct {
	character string
	meanings  []string
	kun       []string
	on        []string
	level     string
}

var all_kanji []kanji

func main() {
	for level := 1; level <= 5; level++ {
		scrapeJisho(level)
	}

	fmt.Println(all_kanji, len(all_kanji))
}

func makeURL(level int) (string, string) {
	_level := "n" + strconv.Itoa(level)
	url := "https://jisho.org/search/%23kanji%20"
	return url + _level, _level
}

func scrapeJisho(level int) {
	c := colly.NewCollector()

	c.OnRequest(func(req *colly.Request) {
		fmt.Println("Visiting", req.URL.String())
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	urlToVisit, kanjiLevel := makeURL(level)

	c.OnHTML("div.kanji_light", func(e *colly.HTMLElement) {
		kanji_char := e.ChildText("div.literal_block")
		kanji_meanings := strings.Split(e.ChildText("div.meanings"), ",")
		kanji_kun := strings.Split(e.ChildText("div.kun"), "、")
		kanji_on := strings.Split(e.ChildText("div.on"), "、")
		kanji_level := kanjiLevel

		kanji_obj := kanji{
			character: kanji_char,
			meanings:  kanji_meanings,
			kun:       kanji_kun,
			on:        kanji_on,
			level:     kanji_level,
		}
		all_kanji = append(all_kanji, kanji_obj)
	})

	c.OnHTML("a.more", func(e *colly.HTMLElement) {
		nextPage := e.Request.AbsoluteURL(e.Attr("href"))
		c.Visit(nextPage)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Response:")
	})

	c.Visit(urlToVisit)
}
