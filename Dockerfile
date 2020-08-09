FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/mikimowski/TWFjaWVqLU1pa3XFgmE

WORKDIR /go/src/github.com/mikimowski/TWFjaWVqLU1pa3XFgmE

RUN go get -d -v ./...
RUN go install -v ./...
#RUN go install

ENTRYPOINT /go/bin/TWFjaWVqLU1pa3XFgmE debug
EXPOSE 8080