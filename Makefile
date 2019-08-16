graphviz:
	sudo apt install graphviz

cpu:
	go tool pprof --pdf cpu.prof > cpu.pdf

mem:
	go tool pprof --pdf -alloc_space mem.prof > mem.pdf

trace:
	go test -test.run TestConvolve
	go tool trace trace.out

bench:
	go test -cpuprofile cpu.prof -memprofile mem.prof -bench .

ui-cpu:
	go tool pprof -http=:8080 cpu.prof

ui-mem:
	go tool pprof -http=:8080 mem.prof

docker-build:
	docker build -t test/golina:latest -f Dockerfile .

docker-run:
	docker run -it --rm \
	-v $(pwd):/wip \
	--entrypoint "/bin/sh" \
	-m 500M \
	--cpus="1" \
	test/golina:latest

gctrace:
	GODEBUG=gctrace=1 go run matrix_test.go 2> stderr.log
