package main

func main() {

	a := App{}
	a.Initialize(
		"postgres",
		"testpass",
		"whattodo")

	a.Run(":9090")

}
