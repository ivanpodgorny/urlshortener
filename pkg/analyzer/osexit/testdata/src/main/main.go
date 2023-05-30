package main

import "os"

func main() {
	os.Exit(0) // want "call of os.Exit in main"
}

//lint:ignore U1000 тестовый файл
//goland:noinspection GoUnusedFunction
func notMain() {
	os.Exit(0)
}
