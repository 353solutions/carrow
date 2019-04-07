.PHONY: clean lib python-bindings

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

python-bindings:
	cd ./python-bindings && \
	go build -o carrow_bindings.so -buildmode=c-shared carrow_bindings.go && \
	pip install -r requirements.txt && \
	python3.6 setup.py build_ext --inplace

python-bindings-clean:
	cd python-bindings && rm -rf *.so *.cpp *.c *.h build

benchmark:
	go test  -run  Example -count 10000

fresh: clean all
