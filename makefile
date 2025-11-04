build:
	@go build -o terminal-note .

run: build
	@./terminal-note