FROM golang:1.18

WORKDIR /code/server
# this is used for hot-reloading applications on change
# configuration is found in the application directories
# themselves
RUN go install github.com/markbates/refresh@latest
