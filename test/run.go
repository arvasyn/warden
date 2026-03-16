package main

import "github.com/arvasyn/warden/internal/photon"

func main() {
	out := photon.Parse("./test/application/manifest.yml")

	err := photon.Run(out, []string{""}, "./test/application")
	if err != nil {
		println(err.Error())
		return
	}

	println("it passed!")
}
