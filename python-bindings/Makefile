.PHONY:

all:
	go build -o libcarrow_bindings.so -buildmode=c-shared carrow_bindings.go && \
#	pip install -r requirements.txt && \
	python setup.py build_ext --inplace

clean:
	rm -rf *.so *.cpp *.c *.h build

test:
	LD_LIBRARY_PATH=${PWD} && python runner.py
