package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func main() {
	randomString := uuid.New().String()

	for {
		currentTime := time.Now().Format(time.RFC3339Nano)
		fmt.Printf("%s: %s\n", currentTime, randomString)
		time.Sleep(5 * time.Second)
	}
}
