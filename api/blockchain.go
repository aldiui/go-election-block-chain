package api

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"go-election-block-chain/domain"
)

type blockchainAPI struct {
	bc *domain.BlockChain
}

func NewBlockchain(app *fiber.App, bc *domain.BlockChain) {
	bca := &blockchainAPI{
		bc: bc,
	}

	app.Get("/chain", bca.Chain)
	app.Post("/give-mandate", bca.GiveMandate)
	app.Get("/check-mandate", bca.CheckMandate)
}

func (bca blockchainAPI) Chain(ctx *fiber.Ctx) error {
	return ctx.JSON(Response[[]*domain.Block]{
		Message: "success",
		Data:    bca.bc.Chain,
	})
}

func (bca blockchainAPI) GiveMandate(ctx *fiber.Ctx) error {
	var req domain.Mandate
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(Response[any]{
			Message: "invalid request",
			Data:    nil,
		})
	}

	isSuccess := bca.bc.GiveMandate(req.From, req.To, req.Value)
	if isSuccess {
		return ctx.JSON(Response[any]{
			Message: "success give mandate",
			Data:    nil,
		})
	}

	return ctx.Status(fiber.StatusInternalServerError).JSON(Response[any]{
		Message: "insufficient mandate",
		Data:    nil,
	})
}

func (bca blockchainAPI) CheckMandate(ctx *fiber.Ctx) error {
	q := ctx.Query("q")

	data := make(map[string]int64)
	for _, v := range strings.Split(q, ",") {
		data[v] = bca.bc.CalculateMandate(v)
	}

	return ctx.JSON(Response[map[string]int64]{
		Message: "success",
		Data:    data,
	})
}
