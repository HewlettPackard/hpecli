package logging

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {

	goPath := os.Getenv("GOPATH")

	var file, err = os.OpenFile(goPath+"/logs/cfmLogs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Could Not Open Log File : " + err.Error())
	}
	log.SetOutput(file)
	log.SetFormatter(&log.TextFormatter{})
}

// Info ...
func Info(info string) {
	log.Info(info)
}

// Error ...
func Error(err string) {
	log.Error(err)
}

// Fatalf ...
func Fatalf(err string) {
	log.Fatalf(err)
}
