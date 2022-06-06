/*
    copied and modified from https://gist.github.com/stanislaw/f62c36823242c4ffea1b
*/
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "utf8/utf8.h"

// check if the first unicode character in c is a single-byte character
bool utf8_is_single_byte(char* c) {
    return (c[0] & 0x80) == 0x0;
}

// check if the first unicode character in c is a double-byte character
bool utf8_is_double_byte(char* c) {
    return (c[0] & 0xe0) == 0xc0 && utf8_is_continuation(c[1]);
}

// check if the first unicode character in c is a triple-byte character
bool utf8_is_triple_byte(char* c) {
    return (c[0] & 0xf0) == 0xe0 && utf8_is_continuation(c[1]) && utf8_is_continuation(c[2]);
}

// check if the first unicode character in c is a quadruple-byte character
bool utf8_is_quadruple_byte(char* c) {
    return (c[0] & 0xf8) == 0xf0 && utf8_is_continuation(c[1]) && utf8_is_continuation(c[2]) && utf8_is_continuation(c[3]);
}

bool utf8_is_continuation(char c) {
    return (c & 0xc0) == 0x80;
}

// returns the number of unicode characters in s
// s must be null-terminated
size_t utf8_strlen(char* s) {
    size_t i = 0, len = 0;
    while(s[i] != 0) { // while not at null-terminator
        if ( ! utf8_is_continuation(s[i])) ++len;
        ++i;
    }
    return len;
}

// returns the number of bytes of the first unicode character in s
// s must be null-terminated
size_t utf8_num_bytes(char* s) {
    size_t len = strlen(s), num_bytes = 0;
    
    // is valid single byte (ie 0xxx xxxx)
    if (len >= 1 && utf8_is_single_byte(s)) {
        num_bytes = 1;
        
    // or is valid double byte (ie 110x xxxx and continuation byte)
    } else if (len >= 2 && utf8_is_double_byte(s)) {
        num_bytes = 2;
        
    // or is valid tripple byte (ie 1110 xxxx and continuation byte)
    } else if (len >= 3 && utf8_is_triple_byte(s)) {
        num_bytes = 3;
        
    // or is valid tripple byte (ie 1111 0xxx and continuation byte)
    } else if (len >= 4 && utf8_is_quadruple_byte(s)) {
        num_bytes = 4;
    }
    
    return num_bytes;
}

// decodes the unicode character c into s
// s must be at least 5 chars long and will be null-terminated by the functions
// returns the number of bytes in c (the resulting strlen of s)
// returns -1 if c is not a valid utf8 character
int utf8_char_to_string(char* s, int32_t c) {
    // copied from the Go stdlib (see utf8.go in unicode/utf8)
    uint32_t ch = (uint32_t)c;
    int num_bytes = 0;
	if (ch <= 127) {
		s[0] = (char)c;
		num_bytes = 1;
	} else if (ch <= 2047) {
		s[0] = 192 | (char)(c >> 6);
		s[1] = 128 | ((char)c)&63;
		num_bytes = 2;
	} else if (ch > 1114111 || (55296 <= ch && ch <= 57343)) {
		return -1;
	} else if (ch <= 65535) {
		s[0] = 224 | (char)(c >> 12);
		s[1] = 128 | ((char)c >> 6)&63;
		s[2] = 128 | ((char)c)&63;
		num_bytes = 3;
	} else { \
		s[0] = 240 | (char)(c>>18);
		s[1] = 128 | (char)(c >> 12);
		s[2] = 128 | ((char)c >> 6)&63;
		s[3] = 128 | ((char)c)&63;
	}
	s[num_bytes] = '\0';
}