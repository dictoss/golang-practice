package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func main() {
	fd, err := os.Open("./data.csv")
	if err == nil {
		// open success
		defer fd.Close()
	} else {
		panic(err)
	}

	fmt.Println("Start")
	fmt.Println("----------")

	reader := csv.NewReader(fd)
	reader.Comma = ','
	reader.LazyQuotes = true

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else {
			if nil == err {
				// success
				fmt.Println(record[0])
				fmt.Println(record[1])
				fmt.Println(record[2])
				fmt.Println("--------")
			} else {
				fmt.Println("read error.")
			}
		}
	}
}
