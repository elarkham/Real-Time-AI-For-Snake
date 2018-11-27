
GOFLAGS =
LDFLAGS =

all: snake

snake: snake.go
	go build $(GOFLAGS) -o snake

clean:
	rm -f snake

