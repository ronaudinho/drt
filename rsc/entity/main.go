package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/dotabuff/manta"
)

type class struct {
	Count  int
	Sample map[string]interface{}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("id not supplied")
	}
	id := os.Args[1]

	f, err := os.Open(fmt.Sprintf("%s.dem", id))
	if err != nil {
		log.Fatalf("unable to open file: %v", err)
	}
	defer f.Close()

	p, err := manta.NewStreamParser(f)
	if err != nil {
		log.Fatalf("unable to create parser: %v", err)
	}

	classes := make(map[string]class)
	p.OnEntity(func(e *manta.Entity, op manta.EntityOp) error {
		c := e.GetClassName()
		if _, ok := classes[c]; !ok {
			classes[c] = class{1, e.Map()}
		} else {
			tmp := classes[c]
			tmp.Count++
			classes[c] = tmp
		}
		return nil
	})

	p.Start()

	b, _ := json.MarshalIndent(classes, "", "  ")
	os.WriteFile(fmt.Sprintf("%s.json", id), b, 0666)
}
