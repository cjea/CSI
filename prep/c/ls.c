#include "pkg/cols.c"
#include "pkg/formatter.c"
#include "pkg/immutable_str.c"
#include "pkg/read_dir.c"
#include <dirent.h>
#include <errno.h>
#include <grp.h>
#include <pwd.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>

#define NUM_COLUMNS 4

int main(int argc, char *argv[]) {
  char *dir_name = argv[1];
  if (dir_name[strlen(dir_name) - 1] != '/') {
    printf("input must be a directory (end with a slash)\n");
    exit(1);
  }
  DIR *dirp = opendir(dir_name);
  if (errno > 0) {
    printf("exiting early. failed to open dir %s\n", dir_name);
    exit(1);
  }

  struct dirent **dirents = dirents_in_dir(dirp);
  row_builder *fptr[NUM_COLUMNS] = {&full_mode, &username, &groupname,
                                    &filesize};

  for (int i = 0; dirents[i] != NULL; i++) {
    char *basename = dirents[i]->d_name;
    if (strlen(basename) > 0) {
      char *abs_path = strcat_i(dir_name, basename);
      format_and_print_row(NUM_COLUMNS, fptr, abs_path);
      free(abs_path);
    }
  }
  closedir(dirp);
  exit(0);
}
