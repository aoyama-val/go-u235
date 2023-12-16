program = ./u235

.PHONY: all
all: build run

.PHONY: build
build: $(program)

$(program): main.go model/model.go
	go build -o $(program)

.PHONY: clean
clean:
	rm -f $(program)

.PHONY: run
run:
	$(program) $(arg)
