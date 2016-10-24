package parser

import (
	"log"
	"os"
	"regexp"
	"time"
	"strings"
	"bytes"
)

type a struct {
	session   string
	starttime string
	endtime   string
}

func handleSession(name string) {
	filepath := "/msys64/records/" + name
	log.Println("start to handle "+filepath)
	var currentSize int64
    var waitingChars=""
	var sessionStartTime int64
	for {
		

		stat, err := os.Stat(filepath)
		if err != nil {
			log.Println(err)
		}
		if currentSize < stat.Size() {
			//file updated
			f, err := os.Open(filepath)
			defer f.Close()
			if err != nil {
				log.Println(err)
			}
			
			//handle diff
			diff := stat.Size() - currentSize
			b := make([]byte, diff)
			_, err = f.ReadAt(b, currentSize)
			if err != nil {
				log.Println(err)
			}
			//handle the first line for timestamp
			if currentSize==0{
				a := bytes.IndexAny(b,"\n")
				sessionStartTime=getStartTime(string(b[0:a]))
				log.Println(sessionStartTime)
				b = b[a+1:]
			}
            waitingChars,err=do(b,waitingChars)
			currentSize = stat.Size()
		} else {
			//do nothing
		}
		time.Sleep(10 * time.Second)
	}
}

func do(b []byte,waitingChars string) (string,error) {
    var remainChars=""
	log.Println("@@@"+waitingChars+"@@@")
    content:=waitingChars+string(b)
	re := regexp.MustCompile("\\$.*\n")
	allIndex := re.FindAllStringIndex(content, -1)
	log.Println(allIndex)
	if(len(allIndex)==0){
		return content, nil
	}
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

		// output := content[cur[1]:next[0]]
        log.Printf(">>>>%v,  %v\n",strings.TrimSpace(command),len(strings.TrimSpace(command)))
		// startTime := gettime(cur[1])
		// endTime := gettime(next[0])
	}
    return remainChars,nil
}

func gettime(index int) time.Time {
	return time.Now()
}

func getStartTime(line string) int64{
	timeString := line[18:]
	const longForm = "Mon, Jan _2, 2006 3:04:05 PM"
    t, _ := time.Parse(longForm, timeString)
    return t.UnixNano()
}