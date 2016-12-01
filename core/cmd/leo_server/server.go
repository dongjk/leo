package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dongjk/leo/core/pkg/storage"
)

var ds *storage.DataStore

func handleChromeInfo(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	_, _ = r.Body.Read(b)

	ds.Insert("chrome", storage.ChromeInfo{time.Now().UnixNano(), string(b)})
	log.Println(string(b))
}

func main() {
	ds, _ = storage.ConstructDataStore()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	http.HandleFunc("/", handleChromeInfo)   // set router
	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
