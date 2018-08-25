package main

import (
	"fmt"
	"flag"
)

var writeBack bool= flag.Bool("-w", false, "write result to (source) file instead of stdout")

func main(){
	flag.Parse()

	for path := flag.Args {
		fmt.Println(path)
	}
}