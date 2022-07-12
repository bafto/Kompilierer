/*
	defines inbuilt ddp functions to work with io
*/
#include <io.h>
#include "ddptypes.h"
#include "utf8/utf8.h"
#include "memory.h"
#include "debug.h"
#include <math.h>

void inbuilt_Schreibe_Zahl(ddpint p1) {
	printf("%ld", p1);
}

void inbuilt_Schreibe_Kommazahl(ddpfloat p1) {
	if (isinf(p1)){
		printf("Unendlich");
	}
	else if (isnan(p1)) {
		printf("Keine Zahl (NaN)");
	}
	else {
		printf("%.16g", p1);
	}
}

void inbuilt_Schreibe_Boolean(ddpbool p1) {
	printf(p1 ? "wahr" : "falsch");
}

void inbuilt_Schreibe_Buchstabe(ddpchar p1) {
	char temp[5];
	utf8_char_to_string(temp, p1);
	printf("%s", temp);
}

void inbuilt_Schreibe_Text(ddpstring* p1) {
	printf("%s", p1->str);
}

#ifdef _WIN32
HANDLE* get_stdin_handle() {
	static HANDLE stdin_hndl;
	return &stdin_hndl;
}
#endif // _WIN32

// discards all characters in stdin up to and including '\n' or EOF
void flush_stdin() {
	int c;
	while ((c = getchar()) != '\n' && c != EOF);
}

ddpstring* inbuilt_Lese_Zeile() {
	ddpstring* dstr = ALLOCATE(ddpstring, 1); // up here to log the address in debug mode
	dstr->str = NULL;
	dstr->cap = 0;
	DBGLOG("inbuilt_Lese_Zeile: %p", dstr);

#define MAX_INPUT_LENGTH 255
#ifdef _WIN32 // TODO: change to ReadFile for redirected input, pipes, files, etc.

	static CONSOLE_READCONSOLE_CONTROL crc = {
		.nLength = sizeof(CONSOLE_READCONSOLE_CONTROL),
		.nInitialChars = 0,
		.dwCtrlWakeupMask = (1 << '\r') | (1 << '\n'),
		.dwControlKeyState = 0,
	};

	wchar_t buff[MAX_INPUT_LENGTH];
	unsigned long read;
	while (ReadConsoleW(*get_stdin_handle(), buff, MAX_INPUT_LENGTH, &read, &crc) != 0 && read != 0) {
		int size = WideCharToMultiByte(CP_UTF8, 0, buff, read, NULL, 0, NULL, NULL); // get the required buffer size
		if (size == 0) runtime_error(1, "WideCharToMultiByte (1) failed with code %ld", GetLastError());
		dstr->str = reallocate(dstr->str, dstr->cap, dstr->cap + size);
		if (WideCharToMultiByte(CP_UTF8, 0, buff, read, dstr->str + dstr->cap, size, NULL, NULL) == 0) runtime_error(1, "WideCharToMultiByte (2) failed with code %ld", GetLastError());
		dstr->cap += size;
		if (read < MAX_INPUT_LENGTH) break;
	}
	dstr->str[dstr->cap-1] = '\0';

#else

	char buff[MAX_INPUT_LENGTH]; // buffer for input
	while (true) { // loop if the string is longer than MAX_INPUT_LENGTH
		size_t size = 0;
		// read chars until CRLF or EOF or MAX_INPUT_LENGTH
		while (size < MAX_INPUT_LENGTH && (buff[size] = getchar()) != EOF) {
			size++;
			if (!utf8_is_multibyte(buff[size-1]) && (buff[size-1] == '\r' || buff[size-1] == '\n')) break;
		}
		// add the read size to the string buffer
		dstr->str = reallocate(dstr->str, dstr->cap, dstr->cap + size);
		// copy the read data into the string
		memcpy(dstr->str + dstr->cap, buff, size);
		dstr->cap += size;
		if (size < MAX_INPUT_LENGTH) break;
	}
	dstr->str[dstr->cap-1] = '\0';

#endif
#undef MAX_INPUT_LENGTH

	return dstr;
}

ddpchar inbuilt_Lese_Buchstabe() {
#ifdef _WIN32 // TODO: change to ReadFile for redirected input, pipes, files, etc.
	wchar_t buff[2];
	char mbStr[4];
	unsigned long read;
	ReadConsoleW(*get_stdin_handle(), buff, 1, &read, NULL);
	int size = WideCharToMultiByte(CP_UTF8, 0, buff, read, mbStr, sizeof(mbStr), NULL, NULL);
	mbStr[size] = '\0';
	return utf8_string_to_char(mbStr);
#else
	char temp[5];
	temp[0] = getchar();
	int i = utf8_indicated_num_bytes(temp[0]);
	for (int j = 1; j < i; j++) temp[j] = getchar();
	temp[i] = '\0';
	flush_stdin();
	return utf8_string_to_char(temp);
#endif
}