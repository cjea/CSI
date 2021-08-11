#ifndef FILE_DIR
#define FILE_DIR

typedef void (*dir_doer)(struct dirent*);
void fail(char* msg);
int each_dir(char* path, dir_doer fptr);

#endif
