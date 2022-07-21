all:
	docker build -t libcontainer-test .
	docker run --privileged --security-opt=seccomp=unconfined libcontainer-test