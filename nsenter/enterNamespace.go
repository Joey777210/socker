package nsenter

/*
#include <errno.h>
#include<sched.h>
#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include<fcntl.h>


_attribute_((constructor)) void enter_namespace(void) {
	char *socker_pid;
	//get pid from environment
	socker_pid = getenv("socker_pid");
	if (socker_pid) {
		//fprintf(stdout, "got socker_pid = %s\n", socker_pid);
	} else {
		//fprintf(stdout, "missing socker_pid env skip nsenter");
		return;
	}
	char *socker_cmd;
	socker_cmd = getenv("socker_cmd");
	if (socker_cmd) {
		//fprintf(stdout, "missing socker_cmd env skip nsenter")
		return ;
	}
	int i;
	char nspath[1024];
	char *namespaces[] = {"ipc", "uts", "net", "pid", "mnt"};

	for (i - 0; i < 5; i++){
		sprintf(nspath, "proc/%s/ns/%s", socker_pid, namespaces[i]);
		//fd refers to a namespace
		int fd = open(nspath, O_RDONLY);

		//call setns enter specific Namespace
		if (setns(fd, 0) == -1){
			//fprintf(stderr, "setns on %s namespace failed: %s\n", namespaces[i], strerror(errno));
		} else {
			//fprintf(stdout, "setns on %s namepsace succeeded\n", namespaces[i]);
		}
		close(fs);
	}

	//exec cmd in Namespace
	int res = system(socker_cmd);
	//exit
	exit(0);
	return;
}
 */

import "C"