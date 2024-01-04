package main

import (
	"fmt"
	"time"
)

func main() {
	currentTime := time.Now()
	fmt.Println("Current Time:", currentTime)

	truncatedTime := currentTime.Truncate(time.Hour)
	fmt.Println("Truncated Time (to nearest hour):", truncatedTime)
}
