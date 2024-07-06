package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
	"os"
)

var fiberLambda *fiberadapter.FiberLambda

type StatusStruct struct {
	Status  string  `json:"status"`
	Data    string  `json:"data"`
	Version float64 `json:"v"`
}

func status(c *fiber.Ctx) error {
	err := c.Status(200).JSON(StatusStruct{
		"ok",
		"API is running",
		0.1,
	})
	return err
}

func IsLambda() bool {
	if lambdaTaskRoot := os.Getenv("LAMBDA_TASK_ROOT"); lambdaTaskRoot != "" {
		return true
	}
	return false
}

func main() {
	app := fiber.New()

	app.Get("/", status)
	app.Get("/users", status)

	if IsLambda() {
		fiberLambda = fiberadapter.New(app)
		lambda.Start(Handler)
	} else {
		app.Listen(":3000")
	}
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return fiberLambda.ProxyWithContext(ctx, request)
}
