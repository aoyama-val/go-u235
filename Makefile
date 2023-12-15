program = ./u235

.PHONY: all
all: build run

.PHONY: build
build: main.go
	go build -o $(program)

.PHONY: clean
clean:
	rm -f $(program)

.PHONY: run
run:
	$(program) $(arg)
