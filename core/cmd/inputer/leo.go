package main

import (
	"os"
	"time"
	"fmt"
	"github.com/dongjk/leo/core/pkg/storage"
)

type keyInfo struct {
	Label string
	Time  int64
}

func main() {
	ds, err := storage.ConstructDataStore()
	if err !=nil{
		fmt.Printf("connect DB error， %v", err)
	}
	
	time := time.Now().UnixNano()

	for _, arg := range os.Args[1:] {
		ds.Insert("keyword", &keyInfo{arg, time})
	}

}
