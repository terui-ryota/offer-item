package main

import (
	"log"
	"os"

	"github.com/terui-ryota/offer-item/cmd"
	"github.com/terui-ryota/offer-item/internal/app/grpcserver"
	_ "go.uber.org/automaxprocs"
)

func main() {
	app, err := grpcserver.InitializeApp()
	if err != nil {
		log.Default().Println("grpcserver.InitializeApp: %w", err)
		os.Exit(1)
	}
	cmd.StartApp(app)
}
