.PHONY : mac windows linux linux_arm all mkdir run web
mac: prepare
	 CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -trimpath -o ./build/T2T_Mac/T2T .

windows: prepare
	go generate
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -trimpath -o ./build/T2T_Windows/T2T.exe .

linux: prepare
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -trimpath -o ./build/T2T_Linux/T2T .

linux_arm: prepare
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -trimpath -o ./build/T2T_Linux_Arm/T2T .

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