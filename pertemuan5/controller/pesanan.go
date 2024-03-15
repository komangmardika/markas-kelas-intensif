package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"pertemuan5/model"
)

func GetPesanan(ctx *fiber.Ctx) error {

	data := []model.Pesanan{
		{
			Id:   "1",
			Name: "Komang",
			Meja: 1,
		},
		{
			Id:   "2",
			Name: "Made",
			Meja: 2,
		},
		{
			Id:   "3",
			Name: "Komang",
			Meja: 3,
		},
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}

	return ctx.Send(marshal)
}
