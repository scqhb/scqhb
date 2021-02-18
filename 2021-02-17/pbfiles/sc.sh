
protoc --go_out=../services/ --go_opt=paths=source_relative --go-grpc_out=requireUnimplementedServers=false:../services/ --go-grpc_opt=paths=source_relative prod.proto

go mod edit -replace=google.golang.org/grpc=github.com/grpc/grpc-go@latest
go mod tidy
go mod vendor
go build -mod=vendor

