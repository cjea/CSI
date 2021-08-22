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
#include <unistd.h>

#define MAX_NUM_COLUMNS 10

void usage(int ec) {
  printf(
      "Usage: ls [ -s ] [ -m ] [ -u ] [ -g ] dir\n-s\tfilesize\n-m\tpermission "
      "mode\n-u\tuser owner\n-g\tgroup owner\n");
  exit(ec);
}

void parse_opts(int argc, char *argv[], row_builder **fptr, int *num_columns) {
  int c;
  while ((c = getopt(argc, argv, "msugh")) != -1) {
    switch (c) {
    case 'h':
      usage(0);
      break;
    case 'm':
      fptr[(*num_columns)++] = &full_mode;
      break;
    case 's':
      fptr[(*num_columns)++] = &filesize;
      break;
    case 'u':
      fptr[(*num_columns)++] = &username;
      break;
    case 'g':
      fptr[(*num_columns)++] = &groupname;
      break;
    }
  }
  return;
}

int main(int argc, char *argv[]) {
  if (argc < 2) {
    printf("must specify a directory\n");
    usage(1);
  }

  char *dir_name = argv[argc - 1];
  if (dir_name[strlen(dir_name) - 1] != '/') {
    printf("input directory must end with a slash\n");
    exit(1);
  }

  row_builder *fptr[MAX_NUM_COLUMNS] = {};
  int num_columns = 0;
  parse_opts(argc, argv, fptr, &num_columns);

  DIR *dirp = opendir(dir_name);
  if (errno > 0) {
    printf("exiting early. failed to open dir %s\n", dir_name);
    exit(1);
  }

  struct dirent **dirents = read_dirents_from_dir(dirp);

  for (int i = 0; dirents[i] != NULL; i++) {
    char *basename = dirents[i]->d_name;
    if (strlen(basename) > 0) {
      char *abs_path = strcat_i(dir_name, basename);
      format_and_print_row(num_columns, fptr, abs_path);
      free(abs_path);
    }
  }
  closedir(dirp);
  exit(0);
}
