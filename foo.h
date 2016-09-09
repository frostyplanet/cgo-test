#ifndef FOO_H
#define FOO_H

#include <pthread.h>
#include <stdint.h>

struct foo {
  pthread_t ths[100];
  int fds[2];
  pthread_mutex_t mutex;
};

typedef struct foo foo_t;

int
foo_init(foo_t *foo);

void
foo_work();

#endif
