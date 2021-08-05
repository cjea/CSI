#include <stddef.h>
#include <errno.h>
#include <stdio.h>
#include <dirent.h>

int main() {
	DIR *dirp;
	struct dirent *dp;

	dirp = opendir("/Users/cjapel/misc/csi/prep/lsafe");

	while (dirp) {
		errno = 0;
		dp = readdir(dirp);
		if (errno == 0) {
			printf("%s", dp->d_name);
			continue;
		}
		if (dp == NULL) {
			printf("DONE\n");
			closedir(dirp);
			break;
		}
		if (errno != 0) {
			printf("Something went wrong: %d\n", errno);
			closedir(dirp);
			return -1;
		}
	}
	return 0;
}
