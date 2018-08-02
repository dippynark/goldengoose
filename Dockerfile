FROM golang:1.6.0-alpine

# copy source code into the $GOPATH and switch to that directory
COPY . ${GOPATH}/src/github.com/dippynark/goldengoose
WORKDIR ${GOPATH}/src/github.com/dippynark/goldengoose

# compile source code and copy into $PATH
RUN go install

# the default command runs the service in the foreground
CMD ["goldengoose"]

# Expose HTTP port and set necessary environment variables
EXPOSE 80