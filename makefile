.PHONY: default build exec

BIN = drawfont
#SRC = $(wildcard *.go)
SRC = main.go

default: build exec

build: $(SRC)
	go build -o $(BIN) $<

exec:
	./$(BIN)
