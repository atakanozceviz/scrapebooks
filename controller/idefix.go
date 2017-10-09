package controller

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/atakanozceviz/scrapebooks/model"
)

type jData struct {
	SearchResponse struct {
		IsFuzzyUsed  bool   `json:"IsFuzzyUsed"`
		IsWordSearch bool   `json:"IsWordSearch"`
		SearchTerm   string `json:"SearchTerm"`
		Products     []struct {
			Name      string  `json:"Name"`
			URL       string  `json:"Url"`
			ImageURL  string  `json:"ImageUrl"`
			Price     float64 `json:"Price"`
			OldPrice  float64 `json:"OldPrice"`
			Discount  float64 `json:"Discount"`
			StockCode string  `json:"StockCode"`
			ID        int     `json:"Id"`
			Persons   struct {
				Person []struct {
					ID          int         `json:"Id"`
					Name        string      `json:"Name"`
					GroupCode   int         `json:"GroupCode"`
					Group       string      `json:"Group"`
					IsDefault   bool        `json:"IsDefault"`
					Description interface{} `json:"Description"`
					ImageURL    interface{} `json:"ImageUrl"`
				} `json:"Person"`
			} `json:"Persons"`
			Categories struct {
				Category []struct {
					ID     int    `json:"Id"`
					Level  int    `json:"Level"`
					Name   string `json:"Name"`
					Parent int    `json:"Parent"`
				} `json:"Category"`
			} `json:"Categories"`
			MediaType            string      `json:"MediaType"`
			MediaTypeText        string      `json:"MediaTypeText"`
			VariationID          string      `json:"VariationId"`
			SoldCount            int         `json:"SoldCount"`
			CommentCount         int         `json:"CommentCount"`
			CreateDate           string      `json:"CreateDate"`
			ComingDate           string      `json:"ComingDate"`
			BrandName            interface{} `json:"BrandName"`
			BrandID              interface{} `json:"BrandId"`
			Rating               int         `json:"Rating"`
			LanguageCode         interface{} `json:"LanguageCode"`
			FirstShipmentDateUtc string      `json:"FirstShipmentDateUtc"`
			StatusCode           int         `json:"StatusCode"`
			ProductProductGroups struct {
				ProductProductGroup []interface{} `json:"ProductProductGroup"`
			} `json:"ProductProductGroups"`
			ManufacturerName string `json:"ManufacturerName"`
			ManufacturerID   int    `json:"ManufacturerId"`
		} `json:"Products"`
		HitCount   int `json:"HitCount"`
		Categories []struct {
			Level           int         `json:"Level"`
			ParentID        int         `json:"ParentId"`
			Seo             string      `json:"Seo"`
			IsBottom        bool        `json:"IsBottom"`
			ParentPath      string      `json:"ParentPath"`
			MetaTitle       interface{} `json:"MetaTitle"`
			MetaDescription string      `json:"MetaDescription"`
			MetaKeywords    interface{} `json:"MetaKeywords"`
			Name            string      `json:"Name"`
			ID              int         `json:"Id"`
			DocumentCount   int         `json:"DocumentCount"`
		} `json:"Categories"`
		Brands []struct {
			Name          string `json:"Name"`
			ID            int    `json:"Id"`
			DocumentCount int    `json:"DocumentCount"`
		} `json:"Brands"`
		MediaTypes []struct {
			Name          string `json:"Name"`
			ID            int    `json:"Id"`
			DocumentCount int    `json:"DocumentCount"`
		} `json:"MediaTypes"`
		Prices []struct {
			From          float64 `json:"From"`
			To            int     `json:"To"`
			Name          string  `json:"Name"`
			ID            int     `json:"Id"`
			DocumentCount int     `json:"DocumentCount"`
		} `json:"Prices"`
		Writers        interface{}   `json:"Writers"`
		PropertyValues []interface{} `json:"PropertyValues"`
		Languages      []interface{} `json:"Languages"`
	} `json:"searchResponse"`
}

func Idefix(books *model.Books, s string) {
	defer wg.Done()
	uri := `http://www.idefix.com/ApiCatalog/Search/`

	jsonStr := []byte(`{"content":{"maxPrice":-1.0,"minPrice":-1.0,"page":1,"size":999,"sortfield":"relevance","sortorder":"desc","term":"` + s + `"},"reponseType":0}`)

	client := &http.Client{}
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(jsonStr))
	LogErr(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "okhttp/2.5.0")

	resp, err := client.Do(req)
	LogErr(err)

	defer resp.Body.Close()
	jsonData, err := ioutil.ReadAll(resp.Body)
	LogErr(err)

	var str jData
	LogErr(json.Unmarshal(jsonData, &str))
	for _, v := range str.SearchResponse.Products {
		for _, c := range v.Categories.Category {
			if c.ID == 11693 {
				title := v.Name
				var author string
				for _, p := range v.Persons.Person {
					if p.GroupCode == 1 {
						author = p.Name
					}
				}

				pub := v.ManufacturerName
				img := "http://i.dr.com.tr" + v.ImageURL
				priceFloat := v.Price
				website := "http://www.idefix.com" + v.URL
				p := model.Book{
					Title:      title,
					Author:     author,
					Publisher:  pub,
					Img:        img,
					Price:      strconv.FormatFloat(priceFloat, 'f', 6, 64),
					PriceFloat: priceFloat,
					WebSite:    website,
					Resource:   "Idefix",
				}
				model.Add(&p, books)
			}
		}

	}

}

func LogErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

//
//func Idefix(books *model.Books, s string) {
//	defer wg.Done()
//	s = strings.Replace(s, " ", "%20", -1)
//	bow := surf.NewBrowser()
//	err := bow.Open("http://www.idefix.com/search?q=" + s)
//	if err != nil {
//		log.Println(err)
//	} else {
//		bow.Find(".list-cell").Each(func(index int, item *goquery.Selection) {
//			a := item.Find(".item-name")
//			title := a.Find("h3").Text()
//			author := item.Find(".who").First().Text()
//			pub := item.Find(".mb10").Text()
//			img, _ := item.Find("figure img").Attr("src")
//			price := item.Find(".price").Text()
//			website, _ := a.Attr("href")
//			if title != "" && price != "" {
//				p := model.Book{
//					Title:     title,
//					Author:    author,
//					Publisher: pub,
//					Img:       img,
//					Price:     price,
//					WebSite:   "http://www.idefix.com" + website,
//					Resource:  "Idefix",
//				}
//				model.Add(&p, books)
//			}
//
//		})
//	}
//}
