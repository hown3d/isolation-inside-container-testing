all:
	docker build -t libcontainer-test .
	docker run --privileged --cap-add SYS_ADMIN --security-opt=apparmor=unconfined --security-opt=seccomp=unconfined libcontainer-test