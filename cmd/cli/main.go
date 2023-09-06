package main

import (
	"ponto-menos/cmd/cli/client"
	_ "ponto-menos/cmd/cli/config/logger"
)

func main() {
	client.Execute()
}
