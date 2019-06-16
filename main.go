package main


import (
	"log"
	"./webserver"
	"os"
	"github.com/joho/godotenv"
)

func  main()  {
	err := godotenv.Load()
	
	ws := webserver.NewWebServer()
	err =  ws.Run(os.Getenv("ADDR"))
	if err != nil {
		log.Printf("Error when start webserver %v", err)
	}
}