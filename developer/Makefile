build-docker:
	docker build . -t carrow:developer
	docker run \
		-v $(PWD):/src/dev \
		-it --workdir=/src/dev/ \
		carrow:developer