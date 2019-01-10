FROM golang:1.10-alpine

# Expose HTTP port and set necessary environment variables
EXPOSE 8000

# copy source code into the $GOPATH and switch to that directory
COPY . ${GOPATH}/src/github.com/dippynark/goldengoose
WORKDIR ${GOPATH}/src/github.com/dippynark/goldengoose

# compile source code and copy into $PATH
RUN go install

# the default command runs the service in the foreground
CMD ["goldengoose"]
