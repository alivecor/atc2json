package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/alivecor/atc2json/atc2json"
)

func main() {
	atcData, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
		return
	}

	jsonOut, err := atc2json.Convert(atcData)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf(jsonOut)

	return
}
