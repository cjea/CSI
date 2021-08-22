#include <errno.h>
#include <stdio.h>
#include <string.h>
#include <sys/stat.h>

#ifndef MAX_ROW_LENGTH
#define MAX_ROW_LENGTH 200
#endif

#ifndef COLUMN_SEPARATOR
#define COLUMN_SEPARATOR "\t"
#endif

// row_builder is any function that can format a stat into a string and append
// it to the given buffer.
typedef void(row_builder)(struct stat, char *buf);

// path_to_stats turns an absolute file path into file stats.
struct stat path_to_stats(char *path) {
  errno = 0;
  struct stat path_stat;
  stat(path, &path_stat);
  if (errno > 0) {
    printf("exiting early. stats call failed for %s\n", path);
    exit(1);
  }

  return path_stat;
}

// format_and_print_row prints output according to the functions in `fptr`, and
// then prints the directory name with a trailing newline.
void format_and_print_row(int num_cols, row_builder **fptr, char *abs_path) {
  struct stat st = path_to_stats(abs_path);
  row_builder **append_col = fptr;
  char row[MAX_ROW_LENGTH] = "";
  for (int i = 0; i < num_cols; i++) {
    (*(append_col + i))(st, row);
    strcat(row, COLUMN_SEPARATOR);
  }
  printf("%s%s\n", row, abs_path);
}
