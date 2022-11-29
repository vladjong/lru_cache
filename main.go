package main

import (
	"fmt"
	"time"
)

func main() {
	var d time.Duration = -10000
	fmt.Println(d.Microseconds())
}
