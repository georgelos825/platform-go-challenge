package main

import (
	"platform-go-challenge/routes"
)

func main() {
	r := routes.SetupRouter()
	// Optional: Improve security by setting trusted proxies
	// Uncomment this line in production to avoid potential IP spoofing attacks
	// err := r.SetTrustedProxies([]string{"127.0.0.1"})
	// if err != nil {
	//     log.Fatalf("Failed to set trusted proxies: %v", err)
	// }
	r.Run(":8080") // Start API on port 8080
}
