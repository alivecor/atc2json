package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/alivecor/atc2json/atc2json"
)

func main() {
	atcData, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
		return
	}

	jsonOut, err := atc2json.Convert(atcData)

	if err != nil {
		log.Fatalln("P#2D4UOQ: ", err)
	}

	fmt.Printf(jsonOut)

	return
}
