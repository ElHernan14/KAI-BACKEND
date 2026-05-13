package main

import (
	"fmt"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("[KAI-API] ")
	fmt.Println("KAI Backend API")
}
