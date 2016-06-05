# RULE 1: Any containers created by this Makefile should be automatically cleaned up

PWD := $(shell pwd)

clean:
	rm bin/entropy-exported
	docker rmi buildertools/entropy:dev

# Using a Node ecosystem tool for service functional tests
prepare:
	docker run -it --rm -v "$(shell pwd)":/work -w /work node:4.1.2 npm install

reiterate:
	docker-compose build entropy
	docker-compose up -d --no-deps --force-recreate entropy

iterate:
	docker-compose kill
	docker-compose rm
	docker-compose up -d

client:
	docker-compose up -d --no-deps --force-recreate client

stop:
	docker-compose stop

# Using `docker cp` to copy a file out of an image requires three steps:
#  1. Create a container from the target image
#  2. Run `docker cp` to copy the file 
#  3. Remove the temporary container
# The biggest problem with this handshake is the need to maintain references to the
# target container. This is compounded in Makefiles. So, forget about that nonsense.
#   Instead I'm using a volume and a copy command from a self-destructing container.
# This has the nice property of being able to run in a single step and potentially 
# performing more complex copy operations.
build:
	docker build -t buildertools/entropy:dev -f build.df .
	docker run --rm -v $(PWD)/bin:/xfer buildertools/entropy:dev /bin/sh -c 'cp /go/bin/entropy* /xfer/'

release: build
	docker build -t buildertools/entropy:latest -f release.df .
	docker tag buildertools/entropy:latest buildertools/entropy:poc

