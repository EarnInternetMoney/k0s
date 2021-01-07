.PHONY: build

NAME     := k0s
PACKAGE  := github.com/btwiuse/$(NAME)
GIT      := $(shell git rev-parse --short HEAD)
DATE     := $(shell date +%FT%T%Z)
VERSION  := $(shell cat VERSION)
IMG_NAME := btwiuse/k0s
IMAGE    := ${IMG_NAME}
LDFLAGS  := $(shell ./pkg/version/ldflags)
TAGS     := ""
SHELL    := bash
BAZEL    := ./bazelw

default: help

fonts:
	@ mkdir -p fonts/ pkg/fonts/; cp /usr/share/figlet/fonts/standard.flf fonts/
	@ assets -d fonts/ -package fonts -o ./pkg/fonts/standard.go -map Fonts
	@ rm -r fonts/

raze:             ## auto generate BUILD.bazel files from Cargo.toml
	@ cargo vendor --versioned-dirs cargo/vendor
	@ which cargo-raze || echo please cargo install cargo-raze | grep --color . 1>&2 && false
	@ cargo raze
	@ git status cargo

gazelle:             ## auto generate BUILD.bazel files from go.mod
	@ go mod tidy
	@ go mod vendor
	@ $(BAZEL) run //:gazelle -- update-repos --from_file=go.mod
	@ $(BAZEL) run //:gazelle
	# @ gazelle -exclude=vendor
	# @ git checkout vendor/nhooyr.io/websocket
	@ cd vendor/google.golang.org/grpc && git grep -e 'importmap = "' | cut -d : -f -1 | sort -u | xargs -L1 -I@ sed -i @ -e 's,importmap = "k0s.io/k0s/vendor/,importmap = ",g'
	@ git status vendor/

bazel-build-android:            ## Build android binaries using bazel
	@ $(BAZEL) build --platforms=@io_bazel_rules_go//go/toolchain:android_amd64  //:k0s
	@ $(BAZEL) build --platforms=@io_bazel_rules_go//go/toolchain:android_386    //:k0s
	@ $(BAZEL) build --platforms=@io_bazel_rules_go//go/toolchain:android_arm    //:k0s
	@ $(BAZEL) build --platforms=@io_bazel_rules_go//go/toolchain:android_arm64  //:k0s

bazel-build-windows:            ## Build windows binaries using bazel
	@ $(BAZEL) build --platforms=@io_bazel_rules_go//go/toolchain:windows_amd64  //:k0s
	@ $(BAZEL) build --platforms=@io_bazel_rules_go//go/toolchain:windows_386    //:k0s

bazel-build-darwin:            ## Build darwin binaries using bazel
	@ $(BAZEL) build --platforms=@io_bazel_rules_go//go/toolchain:darwin_amd64  //:k0s

bazel-build-bsd:            ## Build bsd binaries using bazel
	@ $(BAZEL) build --platforms=@io_bazel_rules_go//go/toolchain:freebsd_amd64  //:k0s
	@ $(BAZEL) build --platforms=@io_bazel_rules_go//go/toolchain:freebsd_386    //:k0s
	@ $(BAZEL) build --platforms=@io_bazel_rules_go//go/toolchain:openbsd_amd64  //:k0s
	@ $(BAZEL) build --platforms=@io_bazel_rules_go//go/toolchain:openbsd_386    //:k0s

bazel-build-bsd-arm:            ## Build bsd arm binaries using bazel
	@ $(BAZEL) build --platforms=@io_bazel_rules_go//go/toolchain:freebsd_arm    //:k0s
	@ $(BAZEL) build --platforms=@io_bazel_rules_go//go/toolchain:openbsd_arm    //:k0s

bazel-build-linux-arm:          ## Build linux arm binaries using bazel
	@ $(BAZEL) build --platforms=@io_bazel_rules_go//go/toolchain:linux_arm64    //:k0s
	@ $(BAZEL) build --platforms=@io_bazel_rules_go//go/toolchain:linux_arm      //:k0s

bazel-build-linux:          ## Build linux binaries using bazel
	@ $(BAZEL) build //:k0s # //cmd/{hub,client,agent}
	@ $(BAZEL) build --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64    //:k0s
	@ $(BAZEL) build --platforms=@io_bazel_rules_go//go/toolchain:linux_arm64    //:k0s
	@ $(BAZEL) build --platforms=@io_bazel_rules_go//go/toolchain:linux_386      //:k0s
	@ $(BAZEL) build --platforms=@io_bazel_rules_go//go/toolchain:linux_mips     //:k0s
	@ $(BAZEL) build --platforms=@io_bazel_rules_go//go/toolchain:linux_mips64   //:k0s
	@ $(BAZEL) build --platforms=@io_bazel_rules_go//go/toolchain:linux_mipsle   //:k0s
	@ $(BAZEL) build --platforms=@io_bazel_rules_go//go/toolchain:linux_mips64le //:k0s
	@ $(BAZEL) build --platforms=@io_bazel_rules_go//go/toolchain:linux_ppc64    //:k0s
	@ $(BAZEL) build --platforms=@io_bazel_rules_go//go/toolchain:linux_ppc64le  //:k0s
	@ $(BAZEL) build --platforms=@io_bazel_rules_go//go/toolchain:linux_s390x    //:k0s

bazel-build:          ## Build binary for current platform using bazel
	@ $(BAZEL) build # //:k0s # //cmd/{hub,client,agent}

build:          ## Build binary for current platform
	@ go run bin.go -tags "$(TAGS)" -ldflags="${LDFLAGS}"

dry:      ## Build binary for every supported platform
	@ go run bin.go -tags "$(TAGS)" -ldflags="${LDFLAGS}" -dry \
	  linux/{{mips{,64},ppc64}{,le},s390x} \
		{linux,android}/{armv6,armv7,arm64,amd64,386} {darwin,windows}/{386,amd64} \
		freebsd/{armv7,armv6} {freebsd,openbsd}/{386,amd64}

build-all:      ## Build binary for every supported platform
	@ make build-android
	@ make build-linux
	@ make build-linux-arm
	@ make build-bsd
	@ make build-bsd-arm
	@ make build-windows
	@ make build-darwin

build-android:  ## Build android binaries
	@ go run bin.go -tags "$(TAGS)" -ldflags="${LDFLAGS}" \
		android/{armv6,armv7,arm64,amd64,386}

build-linux-arm: ## Build linux arm binaries
	@ go run bin.go -tags "$(TAGS)" -ldflags="${LDFLAGS}" \
		linux/{armv6,armv7,arm64}

build-bsd-arm:  	## Build bsd arm binaries
	@ go run bin.go -tags "$(TAGS)" -ldflags="${LDFLAGS}" \
	  freebsd/{armv7,armv6} # ,arm64

build-bsd:  	## Build bsd binaries
	@ go run bin.go -tags "$(TAGS)" -ldflags="${LDFLAGS}" \
	  {freebsd,openbsd}/{amd64,386}

build-linux:  	## Build linux binaries
	@ go run bin.go -tags "$(TAGS)" -ldflags="${LDFLAGS}" \
	  linux/{amd64,386} linux/{{mips{,64},ppc64}{,le},s390x}

build-windows:  ## Build windows binaries
	@ go run bin.go -tags "$(TAGS)" -ldflags="${LDFLAGS}" \
		windows/{386,amd64}

build-darwin:   ## Build darwin binaries
	@ go run bin.go -tags "$(TAGS)" -ldflags="${LDFLAGS}" \
		darwin/amd64

scratch-build:  ## Build without using existing build cache
	@ go run bin.go -d releases/latest -ldflags="${LDFLAGS}" -- -a

scratch-build-all:      ## Build binary for every supported platform ignoring build cache
	@ go run bin.go -tags "$(TAGS)" -ldflags="${LDFLAGS}" -- -a\
		{linux,android}/{armv6,armv7,arm64,amd64,386} {darwin,windows}/{386,amd64} \
	  linux/{{mips{,64},ppc64}{,le},s390x}

release:        ## Build and upload binaries for all supported platforms
	@ mkdir -p bin/; git -C bin/ init
	@ make build-all
	@ pushd bin; tree -L 1 -H '.' --noreport --charset utf-8 > index.html; popd
	@ .ci/release-latest.sh

dist:           ## Build and make an dist image TODO: android builder image
	@ make build
	@ make build-linux-arm
	@ make build-windows
	@ make build-darwin
	@ pushd bin; test -f index.html && rm index.html; tree -H '.' --noreport --charset utf-8 | sponge index.html; popd
	@ cp VERSION bin/
	@ ../../srvhubcli/bin/mnt https://mntcd.op.milvzn.com/btwiuse/k0s/v0.0.10/ ./bin/

install:        ## install binary to system paths
	install -Dvm755 bin/$(NAME) /usr/bin/$(NAME)
	install -Dvm644 .systemd/k0s-agent@.service /etc/systemd/system/k0s-agent@.service
	install -Dvm644 .systemd/k0s-hub.service /etc/systemd/system/k0s-hub.service
	make systemd

systemd:        ## Show systemd post install actions
	@ echo
	@ echo '# Reload systemd unit files'
	@ echo 'sudo systemctl daemon-reload'
	@ echo
	@ echo '# Run this manually to initialize/restart the agent service'
	@ echo 'sudo systemctl enable k0s-agent@$$USER'
	@ echo 'sudo systemctl start k0s-agent@$$USER'
	@ echo 'sudo systemctl restart k0s-agent@$$USER'
	@ echo
	@ echo '# Run this manually to initialize/restart the hub service'
	@ echo 'sudo systemctl enable k0s-hub'
	@ echo 'sudo systemctl start k0s-hub'
	@ echo 'sudo systemctl restart k0s-hub'
	@ echo
	@ echo '# To stop those services，do the following:'
	@ echo 'sudo systemctl list-units | grep k0s'
	@ echo 'sudo systemctl stop k0s-agent@$$USER'
	@ echo 'sudo systemctl stop k0s-hub'

clean:          ## Clean build artifacts
	rm -r bin

up:             ## Start prometheus and grafana to scrape and display metrics
	.docker-compose/up

buildkite:      ## Generate buildkite pipeline yml definition
	cd .buildkite && ./gen | tee /dev/stderr > pipeline.yml

test-build:           ## Check all build tags will compile
	@ go run bin.go -tags raw        # test tag raw
	@ go run bin.go -tags nhooyr     # test tag nhooyr
	@ go run bin.go -tags gorilla    # test tag gorilla

test:           ## Run all tests
	@go clean --testcache && go test ./...


cover:          ## Run test coverage suite
	@go test ./... --coverprofile=cov.out
	@go tool cover --html=cov.out

img:            ## Build Docker Image
	@docker build --rm -t ${IMAGE} .

bazel-img:      ## Build Docker Image with bazel
	@ ${BAZEL} build //:k0s_image
	@ docker load -i bazel-bin/k0s_image.tar
	@ docker tag bazel:k0s_image ${IMAGE}

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":[^:]*?## "}; {printf "\033[38;5;69m%-30s\033[38;5;38m %s\033[0m\n", $$1, $$2}'
