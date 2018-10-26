rm -rf vendor
mkdir  vendor
govendor add +external
go clean
go build 