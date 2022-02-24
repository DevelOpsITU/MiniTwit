FROM golang:bullseye

MAINTAINER Group-d

#ENV GIN_MODE=release
ENV PORT=8080
EXPOSE $PORT


COPY src /minitwit/src
WORKDIR /minitwit/src
RUN cp /minitwit/src/minitwit.db /tmp/minitwit.db

RUN go build -o ./minitwit
ENTRYPOINT [ "./minitwit" ]