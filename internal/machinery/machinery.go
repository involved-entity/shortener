package machinery

import (
	"log"
	"os"

	"github.com/RichardKnop/machinery/v2"
)

var server *machinery.Server

func Init() {
	server = New()
}

func GetServer() *machinery.Server {
	if server == nil {
		log.Fatal("Machinery server not initialized")
		os.Exit(1)
	}
	return server
}
