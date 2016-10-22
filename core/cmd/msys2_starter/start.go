package main

import (
	"log"
	"os/exec"
    "time"
    "fmt"
)

func main() {
    nameSuffix:=fmt.Sprintf("%d",time.Now().UnixNano())
    log.Printf(nameSuffix)
    scriptFileName:="script_"+nameSuffix
    timeFileName:="time_"+nameSuffix

	cmd := exec.Command("/msys64/msys2","script", "-f","--timing="+timeFileName,scriptFileName)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Waiting for command to finish...")
	err = cmd.Wait()
	log.Printf("Command finished with error: %v", err)
}