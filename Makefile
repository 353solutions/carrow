.PHONY: clean python-bindings

all: libcarrow.a
	go build .

libcarrow.a: carrow.o
	ar r $@ $^

%.o: %.cc
	g++ -Wall -O2 -std=c++11 -o $@ -c $^

clean:
	rm -f *.o *.a
	rm -rf ./lib/artifacts

build-docker:
	docker build . -t carrow:builder
	docker run -v $(PWD):/home/carrow -it --workdir=/home/carrow/ carrow:builder

test:
	go test -v ./...

carrow-lib:
	mkdir -p lib/artifacts
	go build -o ./lib/artifacts/libcarrow.so -buildmode=c-shared lib/carrow_lib.go

python-bindings:
	cd ./python-bindings && python3.6 setup.py build_ext --inplace && mkdir -p artifacts && mv -t artifacts *.so *.cpp build

benchmark:
	go test  -run  Example -count 10000

fresh: clean all
