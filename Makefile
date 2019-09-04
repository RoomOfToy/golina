test:
	go test `go list ./... | grep -v examples` -coverprofile=coverage.out
	go tool cover -func=coverage.out

bench:
	go test -bench . `go list ./... | grep -v examples`

doc:
	go doc -all
