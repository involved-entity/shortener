package main

import "shortener/internal"

func main() {
	internal.Run(internal.MustLoad())
}
