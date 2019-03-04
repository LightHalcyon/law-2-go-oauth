ARG GO_VERSION=1.12

FROM golang:${GO_VERSION} AS builder

ADD . /go/src/github.com/reznov53/law-2-go-oauth

ENV OAUTHURL=https://oauth.infralabs.cs.ui.ac.id
ENV CLIENTID=9c6xS7Z1XQWHzkLxMZHxvs0vmy0zFBUK
ENV CLIENTSECRET=hggpGtRNEMuU7nro4Z2WjODfB0Mdm3bc

RUN go get github.com/gorilla/mux
RUN go install github.com/reznov53/law-2-go-oauth

ENTRYPOINT /go/bin/law-2-go-oauth

EXPOSE 8000