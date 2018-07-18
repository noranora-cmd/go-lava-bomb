IMAGE_TAG=v1alpha1
CONFIG_FILE?=/data/config.json
QUAY_PASS?=biggestsecret

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go-lava-bomb main.go
	docker build -t quay.io/tamarakaufler/lava-bomb:$(IMAGE_TAG) .
	docker login quay.io -u tamarakaufler -p $(QUAY_PASS)
	docker push quay.io/tamarakaufler/lava-bomb:$(IMAGE_TAG)

dev:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go-lava-bomb main.go
	docker build -t quay.io/tamarakaufler/lava-bomb:$(IMAGE_TAG) .

run:
	docker run --name=lava-bomb --rm -v $(PWD):/data quay.io/tamarakaufler/lava-bomb:v1alpha1 -file=$(CONFIG_FILE)