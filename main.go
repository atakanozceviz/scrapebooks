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
	//Makes scraper process in parallel
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	//Gets PORT value
	port := os.Getenv("PORT")
	http.HandleFunc("/", json)
	//Serves noimg image
	http.HandleFunc("/noimage", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./view/noimg.png")
	})
	http.HandleFunc("/jsonp/", jsonp)
	//Serves favicon.ico
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./view/favicon.ico")
	})
	//Serves server time now
	http.HandleFunc("/sync", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(strconv.FormatInt(time.Now().UTC().Unix(), 10)))
	})
	//Starts server
	log.Println("Serving on port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

var re = regexp.MustCompile(`(^ +)|( +$)`)

func json(w http.ResponseWriter, r *http.Request) {
	//Gets key from client
	key := r.FormValue("key")
	//Validates the key, debug is left open for debugging
	if controller.CheckKey(key) || key == "debug" {
		//Gets keyword for scraping
		k := html.EscapeString(r.FormValue("keyword"))
		//Beautifies keyword
		k = re.ReplaceAllString(k, "")
		//Checks if keyword is longer than 3 letters
		if len(k) >= 3 {
			books := model.Books{}
			//Starts search and fills with results
			books = *controller.Search(&books, k)
			//Gets average price from results
			avg := books.GetAvg()
			//If nothing found, appends empty values
			if len(books) == 0 {
				avg = 0.0
				books = append(books, model.Book{"Null", "Null", "Null", "https://scrapebooks.herokuapp.com/noimage", "Null", 0.0, "Null", "Null"})
			}

			//Writes JSON to client
			res := model.Result{
				Books: books,
				Avg:   avg,
			}

			w.Header().Set("Content-Type", "application/json;charset=utf-8")
			w.Write(res.ToJson())
		} else {
			w.Write([]byte("Keyword must be longer than 3 letters!"))
		}
	} else {
		w.Write([]byte("Wrong key!"))
	}

}

func jsonp(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	if controller.CheckKey(key) || key == "debug" {
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
