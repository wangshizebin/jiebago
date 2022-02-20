PROJECT = jiebago

.PHONY: all clean build
all:clean build
	@echo "Done!"

build:
	go build -ldflags "-s -w" -buildmode=c-shared -o $(PROJECT).dll $(PROJECT).go

clean:
	rm -rf $(PROJECT).dll $(PROJECT).h
