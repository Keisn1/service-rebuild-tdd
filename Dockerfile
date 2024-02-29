ARG GO_VERSION
FROM golang:${GO_VERSION}-alpine AS base

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -v -o ./note-taking-app .
CMD ./note-taking-app

MAINTAINER  <kay@kayarch>
