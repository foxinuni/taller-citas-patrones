#include <stdio.h>
#include <stdlib.h>
#include <sys/mman.h>
#include <unistd.h>
#include <sys/wait.h>

#define N 68

pid_t fork_count(int* count) {
    pid_t pid = fork();
    if (pid < 0) {
        perror("fork failed");
    }

    if (pid == 0) {
        (*count)++;
    }

    return pid;
}

int main() {
    // mmap nos permite crear memoria compartida entre procesos
    int* count = (int*) mmap(NULL, sizeof(int), PROT_READ | PROT_WRITE, MAP_SHARED | MAP_ANONYMOUS, -1, 0);
    if (count == MAP_FAILED) {
        perror("mmap failed");
    }

    *count = 0;

    int pid[N] = {0};
    for (int i = 0; i < N; i++)
    {
        pid[i] = fork_count(count);
        if (pid[i] == 0)
        { 
            if (i == (int)((N - 1) / 2))
            {
                for (int j = 0; j < N; j++)
                {
                    pid[j] = fork_count(count);
                    if (pid[j] > 0)
                    {
                        wait(NULL);
                        break;
                    }
                }
            }
            exit(0);
        }
    }

    for (int i = 0; i < N; i++)
    {
        wait(NULL);
    }

    printf("count = %d\n", *count);
}