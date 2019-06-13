package main

import (
	"github.com/splisson/opstic"
)




func main() {
    //log := logrus.New()

	r := opstic.BuildEngine()
	//fmt.Printf("Starting opstic server\n")
	r.Run()
}

