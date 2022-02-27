FROM golang:bullseye

MAINTAINER Group-d

#ENV GIN_MODE=release
ENV PORT=8080
EXPOSE $PORT


COPY ./ /minitwit
WORKDIR /minitwit
RUN cp /minitwit/minitwit.db /tmp/minitwit.db

RUN go mod tidy
RUN go mod download

RUN go build -o ./minitwit
ENTRYPOINT [ "./minitwit" ]