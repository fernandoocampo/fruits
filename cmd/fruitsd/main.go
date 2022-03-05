package main

import "github.com/fernandoocampo/fruits/internal/application"

func main() {
	newInstance := application.NewInstance()
	err := newInstance.Run()
	if err != nil {
		panic(err)
	}
	defer newInstance.Stop()
}
