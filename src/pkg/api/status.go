package api

import (
	"github.com/gofiber/fiber/v2"
)

type StatusStruct struct {
	Status  string  `json:"status"`
	Data    string  `json:"data"`
	Version float64 `json:"v"`
}

func StatusHandler(c *fiber.Ctx) error {
	err := c.Status(200).JSON(StatusStruct{
		"ok",
		"API is running",
		0.1,
	})
	return err
}
