# Isolation testing 
This repository is used to try running isolated processes inside docker without using --privileged flag on docker run

## Unprivileged Libcontainer testing
Everything regarding libcontainer (runc) is in libcontainer folder.

### Conclusion
I was not able to provide a setup to run containers without the --privileged flag.

## Bubblewrap
Bubblewrap says that it needs less privileges than runc. Let's see, if that's usable for our purpose.


### Conclusion
Bubblewrap is not possible to run without --privileged, see https://github.com/containers/bubblewrap/issues/505 and https://github.com/containers/bubblewrap/issues/362
