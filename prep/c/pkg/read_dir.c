#include <dirent.h>
#include <errno.h>
#include <stdlib.h>

struct dirent **read_dirents_from_dir(DIR *dirp) {
  struct dirent *dp;
  static struct dirent *dirents[300];
  int num_dirents = 0;

  while (dirp) {
    if ((dp = readdir(dirp)) != NULL) {
      if (num_dirents >= 300) {
        printf("cant handle more than 300 entries until i learn realloc. \n");
        closedir(dirp);
        exit(1);
      }
      dirents[num_dirents++] = dp;
    } else {
      if (errno != 0) {
        printf("exiting early. failed to read directory.");
        exit(1);
      }
      return dirents;
    }
  }
  return dirents;
}
