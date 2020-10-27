FROM scratch
EXPOSE 8000
ENTRYPOINT ["/goldengoose"]
COPY ./bin/ /
