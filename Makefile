install:
	go install -v

fmt:
	go fmt
	cd ./lib && go fmt

image:
	docker build -t cirocosta/l4 .

.PHONY: install
