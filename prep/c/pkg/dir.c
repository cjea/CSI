#include <dirent.h>
#include <errno.h>

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
