.PHONY: clean py

all: libcarrow.a
	go build .

libcarrow.a: carrow.o
	ar r $@ $^

%.o: %.cc
	g++ -Wall -O2 -std=c++11 -o $@ -c $^

clean: python-bindings-clean
	rm -f *.o *.a

build-docker:
	docker build . -t carrow:builder
	docker run -v $(PWD):/home/carrow -it --workdir=/home/carrow/ carrow:builder

test:
	go test -v ./...

circleci:
	docker build -f Dockerfile.test .

py:
	cd py && make

py-clean:
	cd py && make clean

py-test:
	cd py && make test

benchmark:
	go test  -run  Example -count 10000

fresh: clean all
