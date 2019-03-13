docker run --rm -e GOPATH=$GOPATH -v ${HOME}:${HOME} -w $(pwd) -u $(id -u):$(id -g) stratoscale/swagger:v1.0.22 $@
