# Go parameters

GOCMD=GO111MODULE=on go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test

PKG_LIST := `go list ./... | grep -v "vendor"`
PKG_LIST_VET := `go list ./... | grep -v "vendor"`

all: test build
build:
	rm -rf target/
	mkdir target/
	cp comet/conf/comet.toml target/comet.toml
	cp logic/conf/logic.toml target/logic.toml
	$(GOBUILD) -o target/comet comet/cmd/main.go
	$(GOBUILD) -o target/logic logic/cmd/main.go

test:
	$(GOTEST) -v ./...

clean:
	rm -rf target/

run:
	nohup target/logic -conf=target/logic.toml -logtostderr 2>&1 > target/logic.log &
	nohup target/comet -conf=target/comet.toml -logtostderr 2>&1 > target/comet.log &

vet:
	@go vet ${PKG_LIST_VET}

stop:
	pkill -f target/comet
	pkill -f target/logic

# 编译oa二进制 amd64
quick_build_amd:
	@echo '┌ start	quick build amd64'
	@bash ./script/build/quick_build.sh
	@echo '└ end	quick build amd64'

# 编译oa二进制 amr64
quick_build_arm:
	@echo '┌ start	quick build arm64'
	@bash ./script/build/quick_build.sh arm64
	@echo '└ end	quick build arm64'

build_linux_amd:
	rm -rf target/
	mkdir target/
	cp comet/conf/comet.toml target/comet.toml
	cp logic/conf/logic.toml target/logic.toml
	GOOS=linux GOARCH=amd64 GO111MODULE=on GOPROXY=https://goproxy.cn,direct GOSUMDB="sum.golang.google.cn" go build -v -o target/comet comet/cmd/main.go
	GOOS=linux GOARCH=amd64 GO111MODULE=on GOPROXY=https://goproxy.cn,direct GOSUMDB="sum.golang.google.cn" go build -v -o target/logic logic/cmd/main.go
