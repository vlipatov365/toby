FROM ubuntu:latest
LABEL authors="iviac"

ENTRYPOINT ["top", "-b"]