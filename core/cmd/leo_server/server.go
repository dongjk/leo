package main

import (
    "fmt"
    "net/http"
    "log"
    "io/ioutil"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	_,_ = r.Body.Read(b);
    fmt.Println(r.Method)
    fmt.Println(string(b))  // print form information in server side

    fmt.Fprintf(w, "Hello astaxie!") // send data to client side
}

func main() {
    http.HandleFunc("/", sayhelloName) // set router
    err := http.ListenAndServe(":9090", nil) // set listen port
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}