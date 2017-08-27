package main

import (
	"fmt"
	"log"
	"github.com/eqinox76/RiseAndFallOfEmpires/model"
)

func main() {

	space := model.NewSpace()

	fmt.Printf("%+v\n", space)

	out, err := space.Serialize()
	if err != nil {
		log.Fatalln("Failed to encode address book:", err)
	}

	fmt.Println(len(out))

}
