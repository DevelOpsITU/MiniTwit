.PHONY: build clean run_fresh run

BINARY_NAME=group_d_go_app.out

#all: run_fresh
all: run

build:
	./build_app.sh ${BINARY_NAME}

clean:
	go clean
	rm -f out/${BINARY_NAME}

clean-all: clean
	rm -f out/*

run:
	go run src/minitwit.go

run_fresh:
	fresh -c my_fresh_runner.conf

