#include <stdio.h>
#include <stdlib.h>

#include "pkg/dir.c"
#include "pkg/display.c"


int main(int argc, char *argv[]) {
	char* path;
	if (argc == 1) {
		path = ".";
	} else {
		path = argv[1];
	}
	exit(each_dir(path, &display_dir));
}
