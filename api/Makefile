.PHONY: init build dev lint

GO111MODULE=on
PROJECT=point-hub-496713

init:
	go mod download

build:
	go build -o ./bin/ ./cmd/httpserver/main.go

dev:
	FUNCTION_TARGET=RunHTTPServer \
	go run ./cmd/httpserver/main.go

lint:
	@type golangci-lint > /dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run ./...

gcloud-login:
	gcloud auth application-default set-quota-project ${PROJECT}

# 本番環境 Cloud へデプロイ
deploy:
	make __deploy project=${PROJECT}

__deploy:
	gcloud functions deploy point-hub \
	  --gen2 \
	  --project ${project} \
	  --region asia-northeast1 \
	  --runtime go125 \
	  --source ./ \
	  --entry-point RunHTTPServer \
	  --trigger-http \
	  --timeout=540s \
	  --allow-unauthenticated
