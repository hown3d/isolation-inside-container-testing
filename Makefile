.PHONY: libcontainer
libcontainer:
	docker build -f libcontainer/Dockerfile -t libcontainer-test .
	docker run --cap-add SYS_ADMIN --security-opt=apparmor=unconfined --security-opt=seccomp=unconfined libcontainer-test


.PHONY: bubblewrap
bubblewrap:
	docker build -f bubblewrap/Dockerfile -t bubblewrap-test .
	docker run bubblewrap-test
