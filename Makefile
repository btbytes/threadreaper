threadreaper: main.go
	go build

.PHONY:
install: threadreaper
	go install

.PHONY:
clean:
	rm -f threadreaper

