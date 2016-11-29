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
	"github.com/dongjk/leo/core/pkg/storage"
)

const (
	workDir = "/msys64/records/"
)

type ShellParser struct {
	timeline         []offset
	sessionStartTime int64
	globalCharOffset int
	InfoChan         chan storage.ShellInteractive
}

type offset struct {
	timeOffset int64
	charOffset int
}


func (sp *ShellParser) handleSession(name string) {
	filepath := workDir + name
	log.Println("start to handle " + filepath)
	tf := workDir + "time_" + strings.Split(filepath, "_")[1]
	go sp.handleTimeFile(tf)
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
				sp.sessionStartTime = getStartTime(string(b[0:a]))
				log.Println(sp.sessionStartTime)
				b = b[a+1:]
			}
			waitingChars, err = sp.do(b, waitingChars)
			currentSize = stat.Size()
		} else {
			//do nothing
		}
		time.Sleep(10 * time.Second)
	}
}


func (sp *ShellParser) handleTimeFile(filepath string) {
	for {
		fd, err := ioutil.ReadFile(filepath)
		if err != nil {
			log.Println(err)
		}
		lines := strings.Split(string(fd), "\n")
		sp.timeline = []offset{}
		for _, line := range lines {
			if line != "" {
				sp.timeline = append(sp.timeline, offset{getTimeOffset(line), getCharOffset(line)})
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

func (sp *ShellParser) do(b []byte, waitingChars string) (string, error) {
	var remainChars = ""
	// log.Println("@@@" + waitingChars + "@@@")
	content := waitingChars + string(b)
	re := regexp.MustCompile("\\$ .*\n") //TODO regex?
	allIndex := re.FindAllStringIndex(content, -1)
	// log.Println(allIndex)
	if len(allIndex) == 0 {
		i := strings.Index(content, "$ ") //last output index
		if i != -1 && i != 0 {
			lastEndtime := sp.gettime(i)
			lastoutput := content[0:i]
			sp.InfoChan <- storage.ShellInteractive{"", lastoutput, 0, lastEndtime}
			return content[i:len(content)], nil
		} else {
			return content, nil
		}
	}
	if strings.Index(content, "$ ") != 0 {
		lastEndtime := sp.gettime(allIndex[0][1])
		lastoutput := content[0:allIndex[0][1]]
		sp.InfoChan <- storage.ShellInteractive{"", lastoutput, 0, lastEndtime}
	}
	for i := 0; i < len(allIndex); i++ {
		cur := allIndex[i]
		command := content[cur[0]:cur[1]]
		command = strip(command)
		startTime := sp.gettime(cur[1])
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
				endTime = sp.gettime(cur[1] + x)
				remainChars = remainChars[x:len(remainChars)]
				sp.globalCharOffset += cur[1] + x
			} else {
				sp.globalCharOffset += cur[1]

			}
		} else {
			next = allIndex[i+1]
			output = content[cur[1]:next[0]]
			endTime = sp.gettime(next[0])
		}

		sp.InfoChan <- storage.ShellInteractive{command, output, startTime, endTime}
	}
	return remainChars, nil
}

func (sp *ShellParser) gettime(o int) int64 {
	targetOffset := sp.globalCharOffset + o
	var c int
	var t int64
	for _, x := range sp.timeline {
		c += x.charOffset
		t += x.timeOffset
		if c >= targetOffset {
			return sp.sessionStartTime + t*1000
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

func strip(content string) string {
	content  = strings.Trim(content, "$ ")
	content = strings.TrimSpace(content)
	return content
}