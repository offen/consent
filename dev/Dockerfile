# Copyright 2022 - Offen Authors <hioffen@posteo.de>
# SPDX-License-Identifier: Apache-2.0

FROM golang:1.18

WORKDIR /code/server
# this is used for hot-reloading applications on change
# configuration is found in the application directories
# themselves
RUN go install github.com/markbates/refresh@latest
