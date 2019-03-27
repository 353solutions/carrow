all: carrow.so

carrow.so: carrow.go
	go build -buildmode=c-shared -o $@ $<


clean:
	rm -f carrow.so

fresh: clean all
