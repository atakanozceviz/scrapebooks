package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/atakanozceviz/scrapebooks/model"
)

type hData struct {
	MerchantNames      []string    `json:"merchant_names"`
	NewSite            string      `json:"new_site"`
	OrderStore         string      `json:"order_store"`
	OrderCurrency      string      `json:"order_currency"`
	PageDomain         string      `json:"page_domain"`
	PageLanguage       string      `json:"page_language"`
	PageSiteName       string      `json:"page_site_name"`
	PageSiteRegion     string      `json:"page_site_region"`
	SiteType           string      `json:"site_type"`
	PageType           string      `json:"page_type"`
	PageName           string      `json:"page_name"`
	CategoryPath       string      `json:"category_path"`
	PageTitle          string      `json:"page_title"`
	PageURL            string      `json:"page_url"`
	PageReferringURL   string      `json:"page_referring_url"`
	PageQueryString    []string    `json:"page_query_string"`
	IsCanonical        string      `json:"is_canonical"`
	CanonicalURL       interface{} `json:"canonical_url"`
	ProductPrices      []string    `json:"product_prices"`
	ProductUnitPrices  []string    `json:"product_unit_prices"`
	ProductBrands      []string    `json:"product_brands"`
	ProductBrand       string      `json:"product_brand"`
	ProductSkus        []string    `json:"product_skus"`
	ProductIds         []string    `json:"product_ids"`
	ProductTop5        []string    `json:"product_top_5"`
	ProductNames       []string    `json:"product_names"`
	ProductCategoryIds []string    `json:"product_category_ids"`
	ProductCategories  []string    `json:"product_categories"`
	ShippingType       []string    `json:"shipping_type"`
}

var re = regexp.MustCompile("{\".+}")

func Hepsiburada(books *model.Books, s string) {
	var data = hData{}
	defer wg.Done()
	s = strings.Replace(s, " ", "+", -1)
	resp, err := http.Get("http://www.hepsiburada.com/ara?q=" + s)
	if err != nil {
		log.Println(err)
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		} else {
			json.Unmarshal(re.Find(body), &data)
			bow, err := goquery.NewDocument("http://www.hepsiburada.com/ara?q=" + s)
			if err != nil {
				log.Println(err)
			} else {
				bow.Find(".product").Each(func(i int, item *goquery.Selection) {
					img, _ := item.Find("img").Attr("src")
					link, _ := item.Find("a").First().Attr("href")
					if img != "" && link != "" && i < len(data.ProductNames) {
						b := model.Book{
							Title:     data.ProductNames[i],
							Author:    "",
							Price:     data.ProductPrices[i],
							Publisher: data.ProductBrands[i],
							WebSite:   "http://www.hepsiburada.com" + link,
							Img:       img,
							Resource:  "Hepsiburada",
						}
						//75 = "K" Means its a book category
						cat := data.ProductSkus[i][0]
						if cat == 75 {
							model.Add(&b, books)
						}

					}

				})
			}

		}
	}

}
