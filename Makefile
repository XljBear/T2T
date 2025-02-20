.PHONY : mac windows linux linux_arm all mkdir run web
GIT_COMMIT := $(shell git describe --tags --always)
BUILD_TIME := $(shell date '+%Y-%m-%d %H:%M:%S')
GO_VERSION := $(shell go version | awk '{sub("go version ", "");print}')
T2T_VERSION := v0.0.1
mac: prepare
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w -X 'T2T/config.Version=$(T2T_VERSION)' -X 'T2T/config.CommitID=$(GIT_COMMIT)' -X 'T2T/config.BuildTime=$(BUILD_TIME)' -X 'T2T/config.GoVersion=$(GO_VERSION)'" -trimpath -o ./build/T2T_Mac/T2T .

windows: prepare
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -X 'T2T/config.Version=$(T2T_VERSION)' -X 'T2T/config.CommitID=$(GIT_COMMIT)' -X 'T2T/config.BuildTime=$(BUILD_TIME)' -X 'T2T/config.GoVersion=$(GO_VERSION)'" -trimpath -o ./build/T2T_Windows/T2T.exe .

linux: prepare
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X 'T2T/config.Version=$(T2T_VERSION)' -X 'T2T/config.CommitID=$(GIT_COMMIT)' -X 'T2T/config.BuildTime=$(BUILD_TIME)' -X 'T2T/config.GoVersion=$(GO_VERSION)'" -trimpath -o ./build/T2T_Linux/T2T .

linux_arm: prepare
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-s -w -X 'T2T/config.Version=$(T2T_VERSION)' -X 'T2T/config.CommitID=$(GIT_COMMIT)' -X 'T2T/config.BuildTime=$(BUILD_TIME)' -X 'T2T/config.GoVersion=$(GO_VERSION)'" -trimpath -o ./build/T2T_Linux_Arm/T2T .

run:
	go run .
all: mac windows linux linux_arm

web:
	cd  t2t-frontend && yarn && yarn build
prepare:
	mkdir -p ./build/
	mkdir -p ./build/T2T_Mac
	mkdir -p ./build/T2T_Windows
	mkdir -p ./build/T2T_Linux
	mkdir -p ./build/T2T_Linux_Arm