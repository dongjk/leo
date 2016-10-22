package parser

import (
	"log"
	"os"
	"regexp"
	"time"
)

type a struct {
	session   string
	starttime string
	endtime   string
}

func handleSession(name string) {
	filepath := "/d/msys64/records/" + name
	var currentSize int64
    var waitingChars string
	for {
		stat, err := os.Stat(filepath)
		if err != nil {
			log.Println(err)
		}
		if currentSize < stat.Size() {
			//file updated
			f, err := os.Open(filepath)
			if err != nil {
				log.Println(err)
			}
			diff := stat.Size() - currentSize
			b := make([]byte, diff)
			_, err = f.ReadAt(b, currentSize)
			if err != nil {
				log.Println(err)
			}
            waitingChars,err=get(b,waitingChars)
			currentSize = stat.Size()
		} else {
			//do nothing
		}
		time.Sleep(10 * time.Second)
	}
}

func get(b []byte,waitingChars string) (string,error) {
    var remainChars=""
    content:=string(b)
	re := regexp.MustCompile("^\\$.*\n")
	allIndex := re.FindAllStringIndex(content, -1)
	for i := 0; i < len(allIndex); i++ {
		cur := allIndex[i]
		command := content[cur[0]:cur[1]]

		var next []int
		if i+1 == len(allIndex) {
            //last output
			next = []int{len(content), -1}
            remainChars=content[cur[1]:next[0]]
		} else {
			next = allIndex[i+1]
		}

		output := content[cur[1]:next[0]]
        log.Printf("--%v--\n%v \n",command,output)
		// startTime := gettime(cur[1])
		// endTime := gettime(next[0])
	}
    return remainChars,nil
}

func gettime(index int) time.Time {
	return time.Now()
}
