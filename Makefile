.PHONY: clean

all: libcarrow.a
	go build .

libcarrow.a: carrow.o
	ar r $@ $^

%.o: %.cc
	g++ -O2 -I/opt/arrow/include -std=c++11 -o $@ -c $^

clean:
	rm -f *.o *.a

fresh: clean all
