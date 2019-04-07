.PHONY:

all:
	go build -o carrow_bindings.so -buildmode=c-shared carrow_bindings.go && \
	pip install -r requirements.txt && \
	python3.6 setup.py build_ext --inplace

clean:
	rm -rf *.so *.cpp *.c *.h build

test:
	python3.6 runner.py