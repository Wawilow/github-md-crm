package main

import (
	"context"
	api "github-md-crm/pkg/api"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
)

var fiberLambda *fiberadapter.FiberLambda

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return fiberLambda.ProxyWithContext(ctx, request)
}

func main() {
	app := fiber.New()

	app.Get("/api/status", api.StatusHandler)
	app.Get("/api/set", api.GithubSetToken)
	app.Get("/api/files", api.GithubRepoFiles)
	app.Get("/api/file", api.GithubRepoFile)
	app.Get("/redirect", api.GithubRedirect)
	app.Get("/callback", api.GithubCallback)
	app.Get("/api/rep", api.GithubMyRepos)
	app.Post("/api/upl", api.GithubSendFile)

	if IsLambda() {
		fiberLambda = fiberadapter.New(app)
		lambda.Start(Handler)
	} else {
		app.Listen(":3000")
	}
}
