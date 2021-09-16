VER=0.1
CONTAINER="containers.cisco.com/uforsgre/cisshgo"

.PHONY: all
all: cisshgo

.PHONY: all
all: image

.PHONY: image
image:
	docker rmi --force $(CONTAINER):$(VER)
	$(MAKE) cisshgo
	docker build -t $(CONTAINER):$(VER) .

push-image:
	docker push $(CONTAINER):$(VER)

cisshgo: cissh.go
	go build

.PHONY: update-go-deps
update-go-deps:
	@echo ">> updating Go dependencies"
	@for m in $$(go list -mod=readonly -m -f '{{ if and (not .Indirect) (not .Main)}}{{.Path}}{{end}}' all); do \
		go get $$m; \
	done
	go mod tidy
ifneq (,$(wildcard vendor))
	go mod vendor
endif

.PHONY: clean
clean:
	rm -f cisshgo
