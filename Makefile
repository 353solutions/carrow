.PHONY: clean

all: libcarrow.a
	go build .

libcarrow.a: carrow.o
	ar r $@ $^

%.o: %.cc
	g++ -Wall -O2 -std=c++11 -o $@ -c $^

clean:
	rm -f *.o *.a

get-arrow:
		git clone git://github.com/apache/arrow.git ../arrow
		(cd ../arrow && git checkout apache-arrow-0.13.0)

build-docker:
	docker build . -t carrow:builder
	docker run \
		-v $(PWD):/src/carrow \
		-v $(shell readlink -f ../arrow):/src/arrow \
		-it --workdir=/src/carrow/ \
		carrow:builder

plasma-client:
		g++ plasma.cc \
			$(shell pkg-config --cflags --libs plasma) \
			$(shell pkg-config --cflags --libs arrow) \
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