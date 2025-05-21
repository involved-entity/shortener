package main

import (
	"shortener/internal"
	conf "shortener/internal/config"
)

func main() {
	internal.Run(conf.MustLoad())
}
