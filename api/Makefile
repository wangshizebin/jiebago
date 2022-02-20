PROJECT = jiebago

.PHONY: all clean build docker
all:clean build
	@echo "Done!"

build:
	go build -o $(PROJECT) $(PROJECT).go

clean:
	rm -rf $(PROJECT)
	rm -rf dictionary

docker: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(PROJECT) $(PROJECT).go
	mkdir dictionary
	cp -r ../dictionary/ dictionary/
	docker build . -t -o $(PROJECT):v1
	rm -rf $(PROJECT)
	rm -rf dictionary

