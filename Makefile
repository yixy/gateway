.PHONY: build clean

GO_ENV=CGO_ENABLED=1
GO_MODULE=GO111MODULE=on
VERSION_PKG=github.com/yixy/gateway/version
VERSION=0.0.1
GO_FLAGS=-ldflags="-X ${VERSION_PKG}.Ver=${VERSION} -X '${VERSION_PKG}.Env=`uname -mv`' -X '${VERSION_PKG}.BuildTime=`date`'"
GO=env $(GO_ENV) $(GO_MODULE) go

ifeq ($(GOOS), linux)
	GO_FLAGS=-ldflags="-linkmode external -extldflags -static -X ${VERSION_PKG}.Ver=${VERSION} -X '${VERSION_PKG}.Env=`uname -mv`' -X '${VERSION_PKG}.BuildTime=`date`'"
endif

BUILD_TARGET=target
BUILD_TARGET_PKG_DIR=$(BUILD_TARGET)/gateway-${VERSION}

build:
	# build blade cli
	$(GO) build $(GO_FLAGS) -o $(BUILD_TARGET_PKG_DIR)/gateway .
	cp config.yml $(BUILD_TARGET_PKG_DIR)

# clean all build result
clean:
	$(GO) clean ./...
	rm -rf $(BUILD_TARGET)
