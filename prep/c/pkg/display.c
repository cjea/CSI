#include <dirent.h>
#include <sys/stat.h>
#include <grp.h>
#include <pwd.h>
#include <string.h>
#include "display.h"

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

char* dir_char(mode_t m) {
	char* d = malloc(sizeof(char) * 2);
	if (S_ISDIR(m)) *d = 'd';
	else *d = '-';
	*(d+1) = '\0';
	return d;
}

char* owner_perms(mode_t m) {
	int idx = m & S_IRWXU;
	return perms[idx / 0100];
}

char* g_perms(mode_t m) {
	int idx = m & S_IRWXG;
	return perms[(idx / 010) % 010];
}

char* other_perms(mode_t m) {
	int idx = m & S_IRWXO;
	return perms[idx % 010];
}

char* groupname(struct stat st) {
	struct group *grp = getgrgid(st.st_gid);
	return grp->gr_name;
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
		"%s%s%s%s\t%s\t%s\t%dB\t%s\n",
		dir_char(m), owner_perms(m), g_perms(m), other_perms(m),
		username(st), groupname(st), st.st_size,
		d->d_name
	);
}
