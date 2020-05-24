package port

import (
	"os"
)

//GetPort is
func GetPort() string {
	port := os.Getenv("PORT") //env port
	if port == "" {
		port = "8080" //local Test
	}
	return port
}
