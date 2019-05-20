.PHONY: clean

all: libcarrow.a
	go build .

libcarrow.a: carrow.o
	ar r $@ $^

%.o: %.cc
	g++ -Wall -O2 -std=c++11 -o $@ -c $^

clean:
	rm -f *.o *.a

build-docker:
	docker build . -t carrow:builder
	docker run \
		-v $(PWD):/carrow \
		-v $(shell readlink -f ../arrow):/arrow \
		-it --workdir=/carrow/ \
		carrow:builder

plasma-client:
		g++ plasma.cc \
			$(shell pkg-config --cflags --libs plasma) \
			-I/arrow/cpp/src \
			--std=c++11 \
			-o plasmac

plasma-server:
		plasma_store_server -m 1000000 -s /tmp/plasma&


test:
	go test -v ./...

circleci:
	docker build -f Dockerfile.test .

benchmark:
	go test  -run  Example -count 10000

fresh: clean all