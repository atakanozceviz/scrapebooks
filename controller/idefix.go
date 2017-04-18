package controller

import (
	"log"

	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/atakanozceviz/scrapebooks/model"
	"gopkg.in/headzoo/surf.v1"
)

func Idefix(books *model.Books, s string) {
	defer wg.Done()
	s = strings.Replace(s, " ", "%20", -1)
	bow := surf.NewBrowser()
	err := bow.Open("http://www.idefix.com/search?q=" + s)
	if err != nil {
		log.Println(err)
	} else {
		bow.Find(".list-cell").Each(func(index int, item *goquery.Selection) {
			a := item.Find(".item-name")
			title := a.Find("h3").Text()
			author := item.Find(".who").First().Text()
			pub := item.Find(".mb10").Text()
			img, _ := item.Find("figure img").Attr("src")
			price := item.Find(".price").Text()
			website, _ := a.Attr("href")
			if title != "" && price != "" {
				p := model.Book{
					Title:     title,
					Author:    author,
					Publisher: pub,
					Img:       img,
					Price:     price,
					WebSite:   "http://www.idefix.com" + website,
					Resource:  "Idefix",
				}
				model.Add(p, books)
			}

		})
	}
}
