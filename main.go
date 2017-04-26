package main

import (
	"html"
	"log"
	"net/http"
	"os"
	"regexp"
	"runtime"

	"strconv"
	"time"

	"github.com/atakanozceviz/scrapebooks/controller"
	"github.com/atakanozceviz/scrapebooks/model"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/", json)
	http.HandleFunc("/noimage", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./view/noimg.png")
	})
	http.HandleFunc("/jsonp/", jsonp)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./view/favicon.ico")
	})
	log.Println("Serving on port: " + port)

	http.HandleFunc("/sync", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(strconv.FormatInt(time.Now().UTC().Unix(), 10)))
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

var re = regexp.MustCompile(`(^ +)|( +$)`)

func json(w http.ResponseWriter, r *http.Request) {
	//Gets key from client
	s := r.FormValue("key")
	//Validates the key
	if controller.Secret(s) {
		//Gets keyword for scraping
		k := html.EscapeString(r.FormValue("keyword"))
		//Beautify keyword
		k = re.ReplaceAllString(k, "")
		if len(k) >= 3 {
			books := model.Books{}
			books = *controller.Search(&books, k)
			avg := books.GetAvg()
			if len(books) == 0 {
				avg = 0.0
				books = append(books, model.Book{"Null", "Null", "Null", "https://scrapebooks.herokuapp.com/noimage", "Null", 0.0, "Null", "Null"})
			}
			res := model.Result{
				Books: books,
				Avg:   avg,
			}
			w.Header().Set("Content-Type", "application/json;charset=utf-8")
			w.Write(res.ToJson())
		} else {
			w.Write([]byte("Wrong request!"))
		}
	} else {
		w.Write([]byte("Wrong key!"))
	}

}

func jsonp(w http.ResponseWriter, r *http.Request) {
	s := r.FormValue("key")
	if controller.Secret(s) {
		k := html.EscapeString(r.FormValue("keyword"))
		cb := r.FormValue("callback")
		if k != "" && cb != "" {
			books := model.Books{}
			books = *controller.Search(&books, k)
			avg := books.GetAvg()
			if len(books) == 0 {
				avg = 0.0
				books = append(books, model.Book{"Null", "Null", "Null", "https://scrapebooks.herokuapp.com/noimage", "Null", 0.0, "Null", "Null"})
			}
			res := model.Result{
				Books: books,
				Avg:   avg,
			}
			jp := cb + "(" + string(res.ToJson()) + ")"
			w.Header().Set("Content-Type", "application/json;charset=utf-8")
			w.Write([]byte(jp))
		} else {
			jp := cb + `({"err": "FormEmpty"})`
			w.Write([]byte(jp))
		}
	} else {
		w.Write([]byte("Wrong key!"))
	}
}
