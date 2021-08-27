#include <grp.h>
#include <pwd.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>

char perms[8][4] = {
    {'-', '-', '-', '\0'}, {'-', '-', 'x', '\0'}, {'-', 'w', '-', '\0'},
    {'-', 'w', 'x', '\0'}, {'r', '-', '-', '\0'}, {'r', '-', 'x', '\0'},
    {'r', 'w', '-', '\0'}, {'r', 'w', 'x', '\0'},
};

void dir_char(struct stat st, char *in_buffer) {
  if (S_ISDIR(st.st_mode))
    strcat(in_buffer, "d");
  else
    strcat(in_buffer, "-");
}

void owner_perms(struct stat st, char *in_buffer) {
  int idx = st.st_mode & S_IRWXU;
  strcat(in_buffer, perms[idx / 0100]);
  return;
}

void g_perms(struct stat st, char *in_buffer) {
  int idx = st.st_mode & S_IRWXG;
  strcat(in_buffer, perms[(idx / 010) % 010]);
  return;
}

void other_perms(struct stat st, char *in_buffer) {
  int idx = st.st_mode & S_IRWXO;
  strcat(in_buffer, perms[idx % 010]);
  return;
}

void groupname(struct stat st, char *in_buffer) {
  struct group *grp;
  grp = getgrgid(st.st_gid);
  strcat(in_buffer, grp->gr_name);
  return;
}

void username(struct stat st, char *in_buffer) {
  struct passwd *pws;
  pws = getpwuid(st.st_uid);
  strcat(in_buffer, pws->pw_name);
  return;
}

void filesize(struct stat st, char *in_buffer) {
  char stringified[20] = "";
  snprintf(stringified, 20, "%lldB", st.st_size);
  strcat(in_buffer, stringified);
  return;
}

// full_mode concats the permission string of a file, e.g. "-rwx-rw-rw" onto
// a buffer.
void full_mode(struct stat st, char *in_buffer) {
  dir_char(st, in_buffer);
  owner_perms(st, in_buffer);
  g_perms(st, in_buffer);
  other_perms(st, in_buffer);
  return;
}
