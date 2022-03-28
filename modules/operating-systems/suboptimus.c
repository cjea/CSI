#include <errno.h>
#include <strings.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

#ifndef WORKERS
#define WORKERS 4
#endif

int START = 2, END = 20;
char *TESTS[] = {"brute_force", "brutish", "miller_rabin"};
int num_tests = sizeof(TESTS) / sizeof(char *);

void must_errno(char *prefix) {
  if (errno == 0) {
    return;
  }
  fprintf(stderr, "%s", prefix);
  fprintf(stderr, "%s\n", strerror(errno));
  exit(errno);
}

struct queue {
  int fd_write;
  int fd_read;
};
struct queue_msg {
  int func;
  int body;
  int result;
};

void enqueue(struct queue q, struct queue_msg msg) {
  write(q.fd_write, &msg, sizeof(msg));
}

void dequeue(struct queue q, struct queue_msg *msg) {
  read(q.fd_read, msg, sizeof(*msg));
}

void close_queue(struct queue q) {
  close(q.fd_read);
  close(q.fd_write);
}

struct queue new_queue() {
  struct queue q;
  int queue_fds[2];
  pipe(queue_fds);
  must_errno("create pipe: ");
  q.fd_write = queue_fds[1];
  q.fd_read = queue_fds[0];

  return q;
}

void log_result(struct queue_msg sub) {
 printf(
  "%13s says %d %s prime.\n",
  TESTS[sub.func], sub.body, sub.result ? "is" : "IS NOT"
  );
}

// Goal: Run primality on START..END using a mixture of tests.
//  - Single message queue
//  - Round robin for the tests
int main(int argc, char *argv[]) {
  struct queue work = new_queue();
  struct queue results = new_queue();
  struct queue_msg pub;
  struct queue_msg sub;
  int total_tests = num_tests * (END - START);
  int tests_given = 0;
  int tests_received = 0;

  for (int i = 0; i < WORKERS; i++) {
    switch(fork()) {
      case -1:
        fprintf(stderr, "%s\n", strerror(errno));
        exit(errno);
      case 0:
        dup2(work.fd_read, STDIN_FILENO);
        dup2(results.fd_write, STDOUT_FILENO);
        close_queue(work);
        close_queue(results);
        execl("./primality", "./primality", (char*) 0);
      default:
        pub.body = START + tests_given/num_tests;
        pub.func = (START + tests_given) % num_tests;
        enqueue(work, pub);
        must_errno("enqueue err: ");
        tests_given++;
      }
    }

  close(work.fd_read);
  close(results.fd_write);

  while (tests_received++ < total_tests) {
    dequeue(results, &sub);
    must_errno("dequeue err: ");
    log_result(sub);
    if (tests_given < total_tests) {
      pub.body = START + tests_given/3;
      pub.func = (START + tests_given) % num_tests;
      enqueue(work, pub);
      must_errno("enqueue err: ");
      tests_given++;
    }
  }
  return errno;
}
