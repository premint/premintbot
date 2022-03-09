tidy:
	go mod tidy

run:
	go run main.go

dev:
	make tidy && make run

test:
	go test ./...

build:
	gcloud builds submit --tag gcr.io/premint-343516/premintbot

deploy:
	gcloud run deploy premintbot \
		--image gcr.io/premint-343516/premintbot \
		--platform managed

ship:
	make test && make build && make deploy