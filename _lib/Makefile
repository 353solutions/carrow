ARROW_SRC_DIR := /src/arrow/cpp/src
MINICONDA_ROOT := /opt/miniconda
CXXOPT := -O2

ARROW_INCLUDE = -I$(MINICONDA_ROOT)/include -I$(ARROW_SRC_DIR)

all: libcarrow.a

libcarrow.a: carrow.o
	ar r $@ $^

carrow.o: carrow.cc
	g++ -Wall -g $(CXXOPT) -std=c++11 $(ARROW_INCLUDE) -o $@ -c $^

.PHONY: clean
clean:
	rm -f *.o *.a

.PHONY: fresh
fresh: clean all
