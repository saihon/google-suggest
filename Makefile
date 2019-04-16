NAME = google-suggest
PREFIX = /usr/local/bin
VERSION = $(shell git describe --tags --abbrev=0)
LDFLAGS = -w -s \
	-X 'main.Name=$(NAME)' \
	-X 'main.Version=$(VERSION)'

.PHONY: $(NAME) test clean install uninstall

$(NAME): clean
	go build -ldflags=$(LDFLAGS) -o $(NAME) ./cmd/$(NAME)/main.go

test:
	go test -v ./...

clean:
	$(RM) $(NAME)

install:
	cp -i $(NAME) $(PREFIX)

uninstall:
	$(RM) -i $(PREFIX)/$(NAME)
