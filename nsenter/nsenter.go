package nsenter

/*
#define _GNU_SOURCE
#include <unistd.h>
#include <errno.h>
#include <sched.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>

__attribute__((constructor)) void enter_namespace(void) {
	char *socker_pid;
	socker_pid = getenv("socker_pid");
	if (socker_pid) {
		//fprintf(stdout, "got socker_pid=%s\n", mydocker_pid);
	} else {
		//fprintf(stdout, "missing socker_pid env skip nsenter");
		return;
	}
	char *socker_cmd;
	socker_cmd = getenv("socker_cmd");
	if (socker_cmd) {
		//fprintf(stdout, "got socker_cmd=%s\n", mydocker_cmd);
	} else {
		//fprintf(stdout, "missing socker_cmd env skip nsenter");
		return;
	}
	int i;
	char nspath[1024];
	char *namespaces[] = { "ipc", "uts", "net", "pid", "mnt" };

	for (i=0; i<5; i++) {
		sprintf(nspath, "/proc/%s/ns/%s", socker_pid, namespaces[i]);
		int fd = open(nspath, O_RDONLY);

		if (setns(fd, 0) == -1) {
			//fprintf(stderr, "setns on %s namespace failed: %s\n", namespaces[i], strerror(errno));
		} else {
			//fprintf(stdout, "setns on %s namespace succeeded\n", namespaces[i]);
		}
		close(fd);
	}
	int res = system(socker_cmd);
	exit(0);
	return;
}
 */
import "C"