package main


import (
	"log"
	"./webserver"
)

func  main()  {
	ws := webserver.NewWebServer()
	err :=  ws.Run("8200")
	if err != nil {
		log.Printf("Error when start webserver %v", err)
	}
}