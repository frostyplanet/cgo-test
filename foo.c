#define _GNU_SOURCE
#include <pthread.h>
#include <stdint.h>
#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <fcntl.h>
#include <string.h>
#include <errno.h>

#include <cfuncs.h>
#include "foo.h"

static void
worker(foo_t *foo);

int
foo_init(foo_t *foo) {
  int i, ret;
  ret = pipe2(foo->fds, O_DIRECT);
  if (ret < 0) {
      printf("create pipe error: %s\n", strerror(errno));
      exit(1);
  }
  for (i=0; i<100; i++) {
    ret = pthread_create(foo->ths + i, NULL, (void * (*)(void *))worker, (void *)foo);
    if (ret < 0) {
      printf("cannot create thread\n");
      return -1;
    }
    //pthread_detach(*(pthread_t *)(foo->ths + i));
  }
  return 0;}

void
foo_work() {
  usleep(100);
}

static void
worker(foo_t *foo) {
  void *job_p;
  int res, offset;
  while (1) {
    offset = 0;
    pthread_mutex_lock(&foo->mutex);
    //while (1) {
      res = read(foo->fds[0], &job_p, sizeof(void *));
      //if (res <= 0) {
      //  break;
      //}
      //offset += res;
      //if (offset == sizeof(void *)) {
      //  break;
      //}
    //}
    pthread_mutex_unlock(&foo->mutex);
    if (errno != 0) {
      printf("read from pipe error: %s\n", strerror(errno));
      return;
    }
    if (res == 0) {
      continue;
    }
    if (res < sizeof(void *)) {
      printf("read from pipe only get %d\n", res);
      exit(2);
    }
    foo_work();
    //printf("job_p %lu\n",  job_p);
//    pthread_mutex_lock(&foo->mutex);
    job_done_callback_cgo(job_p);
//    pthread_mutex_unlock(&foo->mutex);
  }
}
