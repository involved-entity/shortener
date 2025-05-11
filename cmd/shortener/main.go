package main

import (
	"shortener/internal"
	conf "shortener/internal/config"

	_ "shortener/docs"
)

//	@title			URL Shortener
//	@version		1.1
//	@description	This is the simple URL Shortener service.

//	@host		localhost:8000
//	@BasePath	/

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				Access JWT Token

// @accept json
// @produce json

func main() {
	internal.Run(conf.MustLoad())
}
