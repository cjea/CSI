#include <stdio.h>
#include <stdlib.h>
#include "leveldb/c.h"

leveldb_t* open_db() {
  leveldb_options_t *options = leveldb_options_create();
  char* name = "/tmp/level_db_test";
  char* err;
  return leveldb_open(options, name, &err);
}

int main() {
  open_db();
  printf("Running!\n");
}
