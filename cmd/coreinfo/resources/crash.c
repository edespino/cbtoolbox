#include <stdio.h>
#include <stdlib.h>

int main() {
    printf("This program will crash and create a core dump.\n");
    int *p = NULL;
    *p = 42; // Dereferencing NULL pointer to induce a crash.
    return 0;
}
