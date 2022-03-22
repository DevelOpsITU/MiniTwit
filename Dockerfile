FROM golang:bullseye

#MAINTAINER Group-d

#ENV GIN_MODE=release
ENV PORT=8080
EXPOSE $PORT


COPY ./ /minitwit
WORKDIR /minitwit


RUN go mod tidy; go mod download; go build -o ./minitwit
ENTRYPOINT [ "./minitwit" ]