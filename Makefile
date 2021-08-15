.PHONY: dist
dist:
	@mkdir -p ./bin
	@rm -f ./bin/*
	GOOS=darwin  GOARCH=amd64 go build -o ./bin/factorygen-darwin64       ./cmd/factorygen
	GOOS=linux   GOARCH=amd64 go build -o ./bin/factorygen-linux64        ./cmd/factorygen
	GOOS=linux   GOARCH=386   go build -o ./bin/factorygen-linux386       ./cmd/factorygen
	GOOS=windows GOARCH=amd64 go build -o ./bin/factorygen-windows64.exe  ./cmd/factorygen
	GOOS=windows GOARCH=386   go build -o ./bin/factorygen-windows386.exe ./cmd/factorygen