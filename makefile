GO = GO111MODULE=on go 
BINARY = ./bin/supermart
MAIN_GO = ./cmd/supermart/main.go


buildonly: 
	@echo "==> Building..."
	@echo "Local build" > version.txt
	@echo `hostname` >> version.txt
	@echo `date` >> version.txt
	cat version.txt
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 ${GO} build -a -o ${BINARY} ${MAIN_GO} 