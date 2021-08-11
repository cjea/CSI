#include <stddef.h>
#include <errno.h>
#include <stdio.h>
#include <dirent.h>
#include <stdlib.h>
#include <sys/types.h>
#include <unistd.h>

#include "pkg/display.c"

typedef void (*dir_doer)(struct dirent*);

void fail(char* msg) {
	printf("%s\n", msg);
	exit(1);
}

int each_dir(char* path, dir_doer fptr) {
	DIR *dirp;
	struct dirent *dp;

	dirp = opendir(path);
	if (errno > 0) {
		fail("bad dir");
	}

	while (dirp) {
		errno = 0;
		dp = readdir(dirp);
		if (errno != 0) {
			closedir(dirp);
			fail("something went wrong");
		} else if (dp)  {
			fptr(dp);
			continue;
		} else {
			closedir(dirp);
			break;
		}
	}
	return 0;
}

int main(int argc, char *argv[]) {
	char* path;
	if (argc == 1) {
		path = ".";
	} else {
		path = argv[1];
	}
	exit(each_dir(path, &display_dir));
}
