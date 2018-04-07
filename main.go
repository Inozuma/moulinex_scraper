package main

import (
	"encoding/json"
	"fmt"
	"log"
)

const (
	moulinexAddr        = "http://www.moulinex.fr"
	moulinexRecipesPath = "/recipe/list/more"
)

func main() {
	q := ":default:product:Cookeo"
	iter := NewRecipeIterator(moulinexAddr, moulinexRecipesPath, q)

	// start recipe number at 1
	n := 1
	for iter.Next() {
		recipe := iter.Recipe()

		data, _ := json.MarshalIndent(recipe, "", "  ")
		fmt.Println(string(data))
		n++
	}

	if err := iter.Err(); err != nil {
		log.Fatal(err)
	}
}
