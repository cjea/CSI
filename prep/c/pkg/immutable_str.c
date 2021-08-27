#include <string.h>

// strcat_i(mmutable) for concatting strings non destructively.
// Caller's responsibility to free the result.
char *strcat_i(const char *s1, const char *s2) {
  size_t size = strlen(s1) + strlen(s2);
  char *ret = (char *)malloc(size + 1);
  snprintf(ret, size + 1, "%s%s", s1, s2);
  return ret;
}
