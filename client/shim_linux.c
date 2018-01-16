#include <stdio.h>
#include <stdlib.h>
#include <stdarg.h>
#include <sys/syscall.h>

#define __USE_GNU
#include <dlfcn.h>


extern int go_inotify_init();
extern int go_inotify_init1(int p0);
extern int go_inotify_add_watch(int p0, char *p1, unsigned int p2);
extern int go_inotify_rm_watch(int p0, int p1);


int inotify_init() {
    return go_inotify_init();
}

int inotify_init1(int p0) {
    return go_inotify_init1(p0);
}

int inotify_add_watch(int p0, char* p1, unsigned int p2) {
    return go_inotify_add_watch(p0, p1, p2);
}

int inotify_rm_watch(int p0, int p1) {
    return go_inotify_rm_watch(p0, p1);
}

long syscall(long number, ...) {
    long (*syscall_real)(long number, ...) = dlsym(RTLD_NEXT, "syscall");

    switch (number) {
    case 253:
        {
            return inotify_init();
        }

    case 254:
        {
            int p0;
            char *p1;
            unsigned int p2;

            va_list args;
            va_start(args, number);
            p0 = va_arg(args, int);
            p1 = va_arg(args, char*);
            p2 = va_arg(args, unsigned int);
            va_end(args);

            return inotify_add_watch(p0, p1, p2);
        }

    case 255:
        {
            int p0, p1;

            va_list args;
            va_start(args, number);
            p0 = va_arg(args, int);
            p1 = va_arg(args, int);
            va_end(args);

            return inotify_rm_watch(p0, p1);
        }

    case 294:
        {
            int p0;

            va_list args;
            va_start(args, number);
            p0 = va_arg(args, int);
            va_end(args);

            return inotify_init1(p0);
        }
    }

    void *arg = __builtin_apply_args();
    void *ret = __builtin_apply((void*)syscall_real, arg, 100);
    __builtin_return(ret);
}
