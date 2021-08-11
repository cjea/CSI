#include <stddef.h>
#include <errno.h>
#include <stdio.h>
#include <dirent.h>
#include <stdlib.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <unistd.h>
#include <pwd.h>
#include <grp.h>

typedef void (*dir_doer)(struct dirent*);

int is_regular_file(struct stat st) {
	return S_ISREG(st.st_mode);
}

int is_dir(mode_t m) {
	return S_ISDIR(m);
}

char* groupname(struct stat st) {
	struct group *grp = getgrgid(st.st_gid);
	return grp->gr_name;
}

char perms[8][4] = {
		{'-', '-', '-', '\0'},
		{'-', '-', 'x', '\0'},
		{'-', 'w', '-', '\0'},
		{'-', 'w', 'x', '\0'},
		{'r', '-', '-', '\0'},
		{'r', '-', 'x', '\0'},
		{'r', 'w', '-', '\0'},
		{'r', 'w', 'x', '\0'},
	};

char dir_char(mode_t m) {
	if (is_dir(m)) return 'd';
	return '-';
}

char* owner_perms(mode_t m) {
	int idx = m & S_IRWXU;
	return perms[idx / 0100];
}

char* other_perms(mode_t m) {
	int idx = m & S_IRWXO;
	return perms[idx % 010];
}

char* g_perms(mode_t m) {
	int idx = m & S_IRWXG;
	return perms[(idx % 0100) / 010];
}

char* username(struct stat st) {
	struct passwd *pws;
	pws = getpwuid(st.st_uid);
	return pws->pw_name;
}

struct stat get_stats(char* path) {
	struct stat path_stat;
	stat(path, &path_stat);
	return path_stat;
}

void display_dir(struct dirent *d) {
	struct stat st = get_stats(d->d_name);
	mode_t m = st.st_mode;
	printf(
		"%c%s%s%s\t%s\t%s\t%dB\t%s\n",
		dir_char(m), owner_perms(m), g_perms(m), other_perms(m), username(st), groupname(st), st.st_size, d->d_name
	);
}

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
