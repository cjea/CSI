#include <errno.h>
#include <math.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>


struct queue_msg {
  int func;  /* index in TESTS */
  int body;
  int result;
};

int brute_force(long n);
int brutish(long n);
int miller_rabin(long n);

int (*TESTS[])(long) = {brute_force, brutish, miller_rabin};

void exit_with_usage() {
  fprintf(stderr, "Usage: ./primality [brute_force|brutish|miller_rabin]\n");
  exit(1);
}

int main(int argc, char*argv[]) {
  struct queue_msg ret;
  long num;
  int (*func)(long), tty;

  tty = isatty(fileno(stdin));
  if (tty) {
    if (argc != 2)
      exit_with_usage();

    if (strcmp(argv[1], "brute_force") == 0)
      func = &brute_force;
    else if (strcmp(argv[1], "brutish") == 0)
      func = &brutish;
    else if (strcmp(argv[1], "miller_rabin") == 0)
      func = &miller_rabin;
    else
      exit_with_usage();

    fprintf(stderr, "Running \"%s\", enter a number:\n> ", argv[1]);

    while (scanf("%ld", &num) == 1) {
      printf("%d\n", (*func)(num));
      fflush(stdout);
      fprintf(stderr, "> ");
    }
  } else {
    for (;;) {
      switch(read(STDIN_FILENO, &ret, sizeof(ret))) {
        case -1:
        	// if (errno == 25) return 0; // Ignore some TTY error when reading past the end of stdin.
          fprintf(stderr, "err: %s\n", strerror(errno));
          exit(errno);
        case 0:
          // printf("read 0 bytes\n");
          exit(errno);
        default:
          func = TESTS[ret.func];
          ret.result = (*func)(ret.body);
          switch(write(STDOUT_FILENO, &ret, sizeof(ret))) {
            case -1:
              exit(errno);
            case 0:
              // printf("Wrote 0 bytes");
              exit(errno);
            default:
            	fflush(stdout);
          };
      }
    }
  }
}

/*
 * Primality test implementations
 */

// Just test every factor
int brute_force(long n) {
  for (long i = 2; i < n; i++)
    if (n % i == 0)
      return 0;
  return 1;
}

// Test factors, up to sqrt(n)
int brutish(long n) {
  long max = floor(sqrt(n));
  for (long i = 2; i <= max; i++)
    if (n % i == 0)
      return 0;
  return 1;
}

int randint(int a, int b) { return rand() % (++b - a) + a; }

int modpow(int a, int d, int m) {
  int c = a;
  for (int i = 1; i < d; i++)
    c = (c * a) % m;
  return c % m;
}

int witness(int a, int s, int d, int n) {
  int x = modpow(a, d, n);
  if (x == 1)
    return 1;
  for (int i = 0; i < s - 1; i++) {
    if (x == n - 1)
      return 1;
    x = modpow(x, 2, n);
  }
  return (x == n - 1);
}

// TODO we should probably make this a parameter!
int MILLER_RABIN_ITERATIONS = 10;

// An implementation of the probabilistic Miller-Rabin test
int miller_rabin(long n) {
  int a, s = 0, d = n - 1;

  if (n == 2)
    return 1;

  if (!(n & 1) || n <= 1)
    return 0;

  while (!(d & 1)) {
    d >>= 1;
    s += 1;
  }
  for (int i = 0; i < MILLER_RABIN_ITERATIONS; i++) {
    a = randint(2, n - 1);
    if (!witness(a, s, d, n))
      return 0;
  }
  return 1;
}
