.PHONY: clean lib python-bindings

all: libcarrow.a lib python-bindings
	go build .

libcarrow.a: carrow.o
	ar r $@ $^

%.o: %.cc
	g++ -Wall -O2 -std=c++11 -o $@ -c $^

clean: lib-clean python-bindings-clean
	rm -f *.o *.a

build-docker:
	docker build . -t carrow:builder
	docker run -v $(PWD):/home/carrow -it --workdir=/home/carrow/ carrow:builder

test:
	go test -v ./...

circleci:
	docker build -f Dockerfile.test .

lib:
	mkdir -p lib/artifacts
	go build -o ./lib/artifacts/libcarrow.so -buildmode=c-shared lib/carrow_lib.go

lib-clean:
	rm -rf ./lib/artifacts

python-bindings:
	cd ./python-bindings && \
	pip install -r requirements.txt && \
	python3.6 setup.py build_ext --inplace && \
	mkdir -p artifacts && mv -t artifacts *.so *.cpp build

python-bindings-clean:
	rm -rf ./python-bindings/artifacts

benchmark:
	go test  -run  Example -count 10000

fresh: clean all
