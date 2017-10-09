package model

import (
	"encoding/json"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Book struct {
	Title      string  `json:"title"`
	Author     string  `json:"author"`
	Publisher  string  `json:"publisher"`
	Img        string  `json:"img"`
	Price      string  `json:"price"`
	PriceFloat float64 `json:"pricefloat"`
	WebSite    string  `json:"website"`
	Resource   string  `json:"resource"`
}

type Books []Book

func (bs Books) Len() int           { return len(bs) }
func (bs Books) Less(i, j int) bool { return bs[i].PriceFloat < bs[j].PriceFloat }
func (bs Books) Swap(i, j int)      { bs[i], bs[j] = bs[j], bs[i] }

type Result struct {
	Books Books   `json:"books"`
	Avg   float64 `json:"avg"`
}

var lock sync.Mutex
var rep = strings.NewReplacer(",", ".", " ", "", "TL", "")

//Adds a book to books
func Add(b *Book, bs *Books) {
	if b.Price != "" && b.PriceFloat == 0.0 {
		pds := rep.Replace(b.Price)
		var err error
		b.PriceFloat, err = strconv.ParseFloat(pds, 64)
		if err != nil {
			log.Println(err)
		}
	}
	lock.Lock()
	*bs = append(*bs, Book{
		Title:      b.Title,
		Author:     b.Author,
		Publisher:  b.Publisher,
		Img:        b.Img,
		Price:      b.Price,
		PriceFloat: b.PriceFloat,
		WebSite:    b.WebSite,
		Resource:   b.Resource,
	})
	lock.Unlock()
}

//Returns JSON value of Result
func (res *Result) ToJson() []byte {
	sort.Sort(res.Books)
	j, err := json.Marshal(*res)
	if err != nil {
		log.Println(err)
	}
	return j
}

//Returns average price of books
func (bs *Books) GetAvg() float64 {
	avg := 0.0
	i := 0.0
	for _, v := range *bs {
		avg += v.PriceFloat
		i++
	}
	return avg / i
}
