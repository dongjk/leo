package parser

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	SessionStorageFile = "session_storage_file.txt"
	SessionList        []string
)

func init() {
	// f, err := os.OpenFile("test.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	// if err != nil {
	// 	fmt.Printf("error opening file: %v", err)
	// }
	// // assign it to the standard logger
	// log.SetOutput(f)
	if _, err := os.Stat(SessionStorageFile); os.IsNotExist(err) {
		// os.Create(SessionStorageFile)
	} else {
		c, err := ioutil.ReadFile(SessionStorageFile)
		if err != nil {
			log.Println(err)
			//Do something
		}
		SessionList = strings.Split(string(c), "\n")
	}
}

func Monitor() {

	for {
		files, _ := filepath.Glob("/msys64/records/script_*")
		for _, f := range files {
			file := filepath.Base(f)
			if containString(SessionList, file) {
				continue
			} else {
				SessionList = append(SessionList, file)
				if err := writeToStorageFile(file); err != nil {
					//do
					log.Println(err)
				}
				go handleSession(file)
			}
		}
		//scan folder every 35 seconds
		time.Sleep(35 * time.Second)
	}
}

func writeToStorageFile(line string) error {
	f, err := os.OpenFile(SessionStorageFile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.WriteString(line + "\n"); err != nil {
		return err
	}
	return nil
}

func containString(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		} else {
			continue
		}
	}
	return false
}
