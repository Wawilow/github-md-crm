export $(grep -v '^#' .env.dev | xargs -0)
cd src
go run .
