package log

import "fmt"

func LogInfo(msg string) {
	fmt.Println("INFO: ", msg)
}

func LogError(err error, msg string) {
	fmt.Println("ERROR: ", err, msg)
}
