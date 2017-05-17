package controller

import (
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/atakanozceviz/scrapebooks/model"
	"gopkg.in/headzoo/surf.v1"
)

func Odakitap(books *model.Books, s string) {
	defer wg.Done()
	s = strings.Replace(s, " ", "+", -1)
	bow := surf.NewBrowser()
	err := bow.Open("https://www.odakitap.com/index.php?p=Products&q=" + s)
	if err != nil {
		log.Println(err)
	} else if _, ok := strconv.ParseFloat(s, 64); ok != nil {
		bow.Find(".liste-item").Each(func(i int, item *goquery.Selection) {
			title := item.Find(".l-name a").Text()
			author := item.Find(".l-owner a").Text()
			pub := item.Find(".l-company a").Text()
			img, _ := item.Find(".liste-photo").Attr("src")
			price := item.Find(".l-price").Text()
			website, _ := item.Find(".liste-photo-wrapper a").Attr("href")

			if title != "" && price != "" {
				p := model.Book{
					Title:      title,
					Author:     author,
					Publisher:  pub,
					Img:        "https://www.odakitap.com" + img,
					Price:      price,
					PriceFloat: 0.0,
					WebSite:    website,
					Resource:   "Odakitap",
				}
				model.Add(&p, books)
			}
		})
	} else {
		item := bow.Find(".main-content")

		title := item.Find(".pd-name").Text()
		author := item.Find(".pd-owner a").Text()
		pub := item.Find(".pd-publisher a span").Text()
		img, _ := item.Find("#main_img").Attr("src")
		price := item.Find("#prd_final_price_display").Text()
		website := bow.Url().String()

		if title != "" && price != "" {
			p := model.Book{
				Title:     title,
				Author:    author,
				Publisher: pub,
				Img:       "https://www.odakitap.com" + img,
				Price:     price,
				WebSite:   website,
				Resource:  "Odakitap",
			}
			model.Add(&p, books)
		}

	}
}
