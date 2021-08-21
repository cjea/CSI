#include <stdio.h>
#include <stdlib.h>
#include <string.h>

char* add_strings(const char* s1, const char* s2) {
	char* ret = (char*) malloc(strlen(s1) + strlen(s2) + 1);
	char* tmp = ret;
	const char* first = s1;
	const char* second = s2;
	while(*first != '\0') {
		*tmp++ = *first++;
	}
	while(*second != '\0') {
		*tmp++ = *second++;
	}
	*tmp = '\0';
	return ret;
}

int main(int argc, char *argv[]) {
	char* path;
	if (argc == 1) {
		path = ".";
	} else {
		path = argv[1];
	}
	char* concatted = add_strings("Hey! ", add_strings("sup, ", "dude!"));
	printf("%s\n", path);
	printf("%s\n", concatted);
}
