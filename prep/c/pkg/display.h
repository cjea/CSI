#ifndef HEADER_DISPLAY
#define HEADER_DISPLAY
#include <sys/stat.h>

struct stat get_stats(char* path);
void display_dir(char *basePath, struct dirent *d);

char* dir_char(mode_t m);
char* owner_perms(mode_t m);
char* other_perms(mode_t m);
char* g_perms(mode_t m);

char* groupname(struct stat st);
char* username(struct stat st);

#endif
