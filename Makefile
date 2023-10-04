.PHONY: default
default: all

.PHONY: all
all: hello

.PHONY: clean
clean:
	make -C backend clean

##############

.PHONY: hello
hello:
	make -C backend hello


