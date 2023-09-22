build:
	go build -o bin/app cmd/awesomeProject/*.go
run: build
	./bin/app
