package controller

import (
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/atakanozceviz/scrapebooks/model"
	"gopkg.in/headzoo/surf.v1"
)

func Pandora(books *model.Books, s string) {
	defer wg.Done()
	bow := surf.NewBrowser()
	url := ""
	if _, err := strconv.ParseInt(s, 10, 64); err == nil {
		url = "http://www.pandora.com.tr/Arama/?type=9&kitapadi=&yazaradi=&yayinevi=&isbn=" + s + "&resimli=1&dil=0&sirala=0"
	} else {
		s = strings.Replace(s, " ", "+", -1)
		url = "http://www.pandora.com.tr/Arama/?type=9&kitapadi=" + s + "&yazaradi=&yayinevi=&isbn=&resimli=1&dil=0&sirala=0"
	}
	err := bow.Open(url)
	if err != nil {
		log.Println(err)
	} else {
		bow.Find(".urunorta").Each(func(index int, item *goquery.Selection) {
			title := item.Find(".kt").Text()
			author := item.Find(".yz").Text()
			pub := item.Find(".yy").Text()
			img, _ := item.Find(".imgcont img").Attr("src")
			price := item.Find(".fyt strong").Text()
			website, _ := item.Find(".imgcont a").Attr("href")
			if title != "" && price != "" {
				p := model.Book{
					Title:      title,
					Author:     author,
					Publisher:  pub,
					Img:        "http://www.pandora.com.tr" + img,
					Price:      price,
					PriceFloat: 0.0,
					WebSite:    "http://www.pandora.com.tr" + website,
					Resource:   "Pandora",
				}
				model.Add(&p, books)
			}
		})
	}
}
