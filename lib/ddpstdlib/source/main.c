/*
	initialization for the ddp-c-runtime
	also defines the entry point of the executable
*/
#include <locale.h>
#include <stdio.h>
#include <stdlib.h>
#include <signal.h>
#include <Windows.h>
#include "gc.h"

// should not be needed in production
// mainly for debugging
static void SegfaultHandler(int signal) {
	printf("Segmentation fault");
	exit(1);
}

// initialize runtime stuff
static void init_runtime() {
	setlocale(LC_ALL, "de_DE"); // enable utf8 printing
	setlocale(LC_NUMERIC, "French_Canada.1252"); // print floats with , instead of . as seperator
#ifdef _WIN32
	// enable utf-8 printing on windows
	// both of the functioncalls below are needed
	SetConsoleCP(CP_UTF8);
	SetConsoleOutputCP(CP_UTF8);
#endif // _WIN32

	signal(SIGSEGV, SegfaultHandler); // "catch" segfaults

	initTable(get_ref_table()); // initialize the reference table
}

// end the runtime
static void end_runtime() {
	freeTable(get_ref_table()); // free the reference talbe (not the remaining entries, it should be empty)
}

extern int inbuilt_ddpmain(); // implicitly defined by the ddp code

// entry point of the final executable (needed by gcc)
int main() {
	init_runtime(); // initialize the runtime
	int ret = inbuilt_ddpmain(); // run the ddp code
	end_runtime(); // end the runtime
	return ret; // return the exit status of the ddp-code
}