package parser

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	workDir = "/msys64/records/"
)

var (
	timeline         []offset
	sessionStartTime int64
	globalCharOffset int
	infoChan         = make(chan info)
)

type offset struct {
	timeOffset int64
	charOffset int
}

type info struct {
	cmd       string
	output    string
	startTime int64
	endtime   int64
}

func handleSession(name string) {
	filepath := workDir + name
	log.Println("start to handle " + filepath)
	tf := workDir + "time_" + strings.Split(filepath, "_")[1]
	go handleTimeFile(tf)
	go printInfo()
	var currentSize int64
	var waitingChars = ""
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
			if currentSize == 0 {
				a := bytes.IndexAny(b, "\n")
				sessionStartTime = getStartTime(string(b[0:a]))
				log.Println(sessionStartTime)
				b = b[a+1:]
			}
			waitingChars, err = do(b, waitingChars)
			currentSize = stat.Size()
		} else {
			//do nothing
		}
		time.Sleep(10 * time.Second)
	}
}

func printInfo() {
	for {
		a := <-infoChan
		if a.cmd != "" {
			log.Printf("CMD==\n%s---%s\n", a.cmd, time.Unix(a.startTime/1000000000, 0))

		}
		if a.output != "" {
			log.Printf("OUTPUT::\n%s---%s\n", a.output, time.Unix(a.endtime/1000000000, 0))

		}
	}
}

func handleTimeFile(filepath string) {
	for {
		fd, err := ioutil.ReadFile(filepath)
		if err != nil {
			log.Println(err)
		}
		lines := strings.Split(string(fd), "\n")
		timeline = []offset{}
		for _, line := range lines {
			if line != "" {
				timeline = append(timeline, offset{getTimeOffset(line), getCharOffset(line)})
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func getTimeOffset(line string) int64 {
	s := strings.Split(line, " ")[0]
	s = strings.Replace(s, ".", "", 1)
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func getCharOffset(line string) int {
	s := strings.Split(line, " ")[1]
	i, _ := strconv.Atoi(s)
	return i
}

func do(b []byte, waitingChars string) (string, error) {
	var remainChars = ""
	// log.Println("@@@" + waitingChars + "@@@")
	content := waitingChars + string(b)
	re := regexp.MustCompile("\\$.*\n") //TODO regex?
	allIndex := re.FindAllStringIndex(content, -1)
	// log.Println(allIndex)
	if len(allIndex) == 0 {
		i := strings.Index(content, "$ ") //last output index
		if i != -1 && i != 0 {
			lastEndtime := gettime(i)
			lastoutput := content[0:i]
			infoChan <- info{"", lastoutput, 0, lastEndtime}
			return content[i:len(content)], nil
		} else {
			return content, nil
		}
	}
	if strings.Index(content, "$ ") != 0 {
		lastEndtime := gettime(allIndex[0][1])
		lastoutput := content[0:allIndex[0][1]]
		infoChan <- info{"", lastoutput, 0, lastEndtime}
	}
	for i := 0; i < len(allIndex); i++ {
		cur := allIndex[i]
		command := content[cur[0]:cur[1]]
		startTime := gettime(cur[1])
		var output string
		var endTime int64
		var next []int
		if i+1 == len(allIndex) {
			//last output
			next = []int{len(content), -1}
			remainChars = content[cur[1]:next[0]]
			x := strings.Index(remainChars, "$ ")
			if x != -1 && x != 0 {
				output = remainChars[0:x]
				endTime = gettime(cur[1] + x)
				remainChars = remainChars[x:len(remainChars)]
				globalCharOffset += cur[1] + x
			} else {
				globalCharOffset += cur[1]

			}
		} else {
			next = allIndex[i+1]
			output = content[cur[1]:next[0]]
			endTime = gettime(next[0])
		}

		infoChan <- info{command, output, startTime, endTime}
	}
	return remainChars, nil
}

func gettime(o int) int64 {
	targetOffset := globalCharOffset + o
	var c int
	var t int64
	for _, x := range timeline {
		c += x.charOffset
		t += x.timeOffset
		if c >= targetOffset {
			return sessionStartTime + t*1000
		}
	}
	return -1
}

func getStartTime(line string) int64 {
	timeString := line[18:]
	const longForm = "Mon, Jan _2, 2006 3:04:05 PM"
	t, _ := time.Parse(longForm, timeString)
	return t.UnixNano()
}
