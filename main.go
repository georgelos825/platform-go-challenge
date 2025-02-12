package main

import (
	"platform-go-challenge/routes"
)

func main() {
	r := routes.SetupRouter()
	r.Run(":8080") // Start API on port 8080
}
