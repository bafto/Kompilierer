#include "ddptypes.h"
#include "ddpwindows.h"
#include "memory.h"
#include <sys/types.h>
#include <dirent.h>
#include <libgen.h>
#include <stdio.h>

// copied from https://stackoverflow.com/questions/11238918/s-isreg-macro-undefined
// to handle missing macros on Windows
#define _CRT_INTERNAL_NONSTDC_NAMES 1
#include <sys/stat.h>
#if !defined(S_ISDIR) && defined(S_IFMT) && defined(S_IFDIR)
  #define S_ISDIR(m) (((m) & S_IFMT) == S_IFDIR)
#endif

#ifdef DDPOS_WINDOWS
#include <io.h>
#include <direct.h>
#define access _access
#define stat _stat
#define mkdir _mkdir
#define F_OK 0
#define PATH_SEPERATOR "/\\"
#else
#include <unistd.h>
#define PATH_SEPERATOR "/"
#define mkdir(arg) mkdir(arg, 0700)
#endif // DDPOS_WINDOWS

ddpbool Existiert_Pfad(ddpstring* Pfad) {
	return access(Pfad->str, F_OK) == 0;
}

ddpbool Erstelle_Ordner(ddpstring* Pfad) {
	// recursively create every directory needed to create the final one
	char* it = Pfad->str;
	while ((it = strpbrk(it, PATH_SEPERATOR)) != NULL) {
		*it = '\0';
		if (mkdir(Pfad->str) != 0) return false;
		*it = '/';
		it++;
	}

	// == '/' because it might have already been created
	return Pfad->str[Pfad->cap - 2] == '/' || mkdir(Pfad->str) == 0;
}

ddpbool Ist_Ordner(ddpstring* Pfad) {
	// remove possible trailing seperators
	char* it = Pfad->str + Pfad->cap - 2; // last character in str
	while (it >= Pfad->str) {
		if (strpbrk(it--, PATH_SEPERATOR) != NULL) {
			*(it+1) = '\0';
		} else {
			break;
		}
	}

	struct stat path_stat;
	if (stat(Pfad->str, &path_stat) != 0) return false;
	return S_ISDIR(path_stat.st_mode);
}

// copied from https://stackoverflow.com/questions/2256945/removing-a-non-empty-directory-programmatically-in-c-or-c
static int remove_directory(const char *path) {
	DIR *d = opendir(path);
	size_t path_len = strlen(path);
	int r = -1;

	if (d) {
		struct dirent *p;
		r = 0;
		while (!r && (p = readdir(d))) {
			int r2 = -1;
			char *buf;
			size_t len;

			// Skip the names "." and ".." as we don't want to recurse on them.
			if (!strcmp(p->d_name, ".") || !strcmp(p->d_name, ".."))
				continue;

			len = path_len + strlen(p->d_name) + 2; 
			buf = malloc(len);

			if (buf) {
				struct stat statbuf;

				snprintf(buf, len, "%s/%s", path, p->d_name);
				if (!stat(buf, &statbuf)) {
					if (S_ISDIR(statbuf.st_mode))
						r2 = remove_directory(buf);
					else
						r2 = unlink(buf);
				}
				free(buf);
			}
			r = r2;
		}
		closedir(d);
	}

	if (!r)
		r = rmdir(path);

	return r;
}

ddpbool Loesche_Pfad(ddpstring* Pfad) {
	if (Ist_Ordner(Pfad)) {
		return remove_directory(Pfad->str) == 0;
	}
	return unlink(Pfad->str) == 0;
}

ddpbool Pfad_Verschieben(ddpstring* Pfad, ddpstring* NeuerName) {
	struct stat path_stat;
	// https://stackoverflow.com/questions/64276902/mv-command-implementation-in-c-not-moving-files-to-different-directory
	if (stat(NeuerName->str, &path_stat) == 0) {
		if (S_ISDIR(path_stat.st_mode)) {
			char* base = basename(Pfad->str);
			Pfad->str = GROW_ARRAY(char, Pfad->str, Pfad->cap, Pfad->cap + strlen(basename) + 1);
			memcpy()
		}
	} else {
		return rename(Pfad->str, NeuerName->str) == 0;
	}
}