.PHONY: phony

build: phony
	gcc -Wall -Wextra -g -o ls.out ls.c

run: build
	./ls.out

go: build
	valgrind -s ./ls.out
