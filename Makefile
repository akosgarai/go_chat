all:
	@echo "Usage: make install"
install:
	go install
	cp index.html $(GOPATH)/index.html
