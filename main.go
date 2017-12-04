package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(time.Now().Unix())


	a := App{}
	a.Initialize(
		"postgres",
		"testpass",
		"whattodo")

	a.Run(":9090")

}
