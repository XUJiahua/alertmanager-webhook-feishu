fmt:
	go fmt ./...
run:fmt
	go run main.go server -c config.yml -e -v
