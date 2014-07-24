ALL := yamlsh

.PHONY: all
.SUFFIXES: .go

all: $(ALL)

.go:
	go build -o $@ $<
