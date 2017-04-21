package controller

import (
	"sync"

	"github.com/atakanozceviz/scrapebooks/model"
)

var wg sync.WaitGroup

func Search(books *model.Books, s string) *model.Books {
	wg.Add(5)
	go Idefix(books, s)
	go Odakitap(books, s)
	go Pandora(books, s)
	go Hepsiburada(books, s)
	go Sozcu(books, s)
	wg.Wait()
	return books
}
