package main

import (
	"html"
	"log"
	"net/http"
	"os"

	"runtime"

	"github.com/atakanozceviz/scrapebooks/controller"
	"github.com/atakanozceviz/scrapebooks/model"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/", json)
	http.HandleFunc("/jsonp/", jsonp)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./view/favicon.ico")
	})
	log.Println("Serving on port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func json(w http.ResponseWriter, r *http.Request) {
	k := html.EscapeString(r.FormValue("keyword"))
	if k != "" {
		var books model.Books
		books = *controller.Search(&books, k)
		avg, err := books.GetAvg()
		if err != nil {
			log.Println(err)
		}
		res := model.Result{
			Books: books,
			Avg:   avg,
		}
		w.Header().Set("Content-Type:", "application/json;charset=utf-8")
		w.Write(res.ToJson())
	} else {
		w.Write([]byte("Wrong request!"))
	}
}

func jsonp(w http.ResponseWriter, r *http.Request) {
	k := html.EscapeString(r.FormValue("keyword"))
	cb := r.FormValue("callback")
	if k != "" && cb != "" {
		var books model.Books
		books = *controller.Search(&books, k)
		//avg, err := books.GetAvg()
		//if err != nil {
		//	log.Println(err)
		//}
		//res := model.Result{
		//	Books: books,
		//	Avg:   avg,
		//}
		jp := cb + "(" + string(books.ToJson()) + ")"
		w.Header().Set("Content-Type:", "application/json;charset=utf-8")
		w.Write([]byte(jp))
	} else {
		w.Write([]byte("Wrong request!"))
	}
}
