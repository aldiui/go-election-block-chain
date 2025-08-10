package main

import (
	"github.com/gofiber/fiber/v2"

	"go-election-block-chain/api"
	"go-election-block-chain/domain"
)

func main() {
	bc := domain.NewBlockChain()

	go bc.PlenaryRecap()

	app := fiber.New()
	api.NewBlockchain(app, bc)

	_ = app.Listen(":8000")
}
