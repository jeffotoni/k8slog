# Makefile
.EXPORT_ALL_VARIABLES:	

AWS_DEFAULT_PROFILE=jeffotoni

GO111MODULE=on
GOPROXY=direct
GOSUMDB=off
GOPRIVATE=github.com/jeffotoni/k8slog

update:
	@echo "########## Compilando nossa API ... "
	@rm -f go.*
	go mod init github.com/jeffotoni/jeffotoni/k8slog
	go mod tidy -compat=1.19 -go=1.19
	CGO_ENABLED=0 GOOS=linux go build --trimpath -ldflags="-s -w"
	@echo "buid update completo..."
	@echo "fim"