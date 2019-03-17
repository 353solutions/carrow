.PHONY: clean

all: libcarrow.a
	go build .

libcarrow.a: carrow.o
	ar r $@ $^

%.o: %.cc
	g++ -I../arrow/cpp/src -O2 -o $@ -c $^

clean:
	rm -f *.o *.a

fresh: clean all
