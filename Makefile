BIN     = conntroll
GOOS    = $(shell go env GOOS)
GOARCH  = $(shell go env GOARCH)
LDFLAGS = $(shell ./pkg/version/ldflags)

all:
	@ go run bin.go -tags "$(TAGS)" -ldflags="${LDFLAGS}"

release:
	@ go run bin.go -d releases/latest -strip -upx -ldflags="${LDFLAGS}" linux/{arm,arm64,amd64,386} darwin/amd64
	@ ./releases/update-latest-index
	@ mkdir -p releases/$(shell git rev-parse HEAD)
	@ cp -rv releases/latest/* releases/$(shell git rev-parse HEAD)
	@ ./releases/update-index
	@ git -C releases add .
	@ git -C releases commit -m $(shell git rev-parse HEAD)

link:
	ln -f bin/$(BIN)-$(GOOS)-$(GOARCH) bin/$(BIN)
	ln -f bin/$(BIN)-$(GOOS)-$(GOARCH) bin/agent
	ln -f bin/$(BIN)-$(GOOS)-$(GOARCH) bin/hub
	ln -f bin/$(BIN)-$(GOOS)-$(GOARCH) bin/client

install:
	install -Dvm755 bin/$(BIN) /usr/bin/$(BIN)
	ln -f /usr/bin/$(BIN) /usr/bin/agent
	install -Dvm644 agent@.service /usr/lib/systemd/system/agent@.service
	make systemd

systemd:
	@ echo 'Run this manually to initialize/restart the agent systemd service'
	@ echo 'sudo systemctl daemon-reload'
	@ echo 'sudo systemctl enable agent@$$USER'
	@ echo 'sudo systemctl start agent@$$USER'
	@ echo 'sudo systemctl restart agent@$$USER'

clean:
	rm -r bin

up:
	.docker-compose/up

buildkite:
	cd .buildkite && ./gen | tee /dev/stderr > pipeline.yml
