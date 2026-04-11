package main

import (
	"github.com/arvasyn/warden/internal/gallium"
	"github.com/rs/zerolog"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	// Gallium is broken and I cba to fix it tonight

	out, err := gallium.Parse("./test/application/manifest.yml")
	if err != nil {
		println(err.Error())
		return
	}

	err = gallium.Run(*out, "./test/application", []string{""})
	if err != nil {
		println(err.Error())
		return
	}

	println("it passed!")
}
