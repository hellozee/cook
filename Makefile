all:
	go build -o cook main.go

.PHONY: clean
clean:
	rm cook

.PHONY: install
install:
	sudo cp cook /usr/bin/

.PHONY: uninstall
uninstall:
	sudo rm /usr/bin/cook
