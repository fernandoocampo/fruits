package main

import "github.com/fernandoocampo/fruits/internal/application"

func main() {
	newInstance := application.NewInstance()
	if err := newInstance.Run(); err != nil {
		panic(err)
	}

	newInstance.Stop()
}
