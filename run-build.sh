#!/bin/bash

go build -buildmode=plugin -o plugins/posts.so gq/posts/*
go build -buildmode=plugin -o plugins/translate.so gq/translate/*
go build -buildmode=plugin -o plugins/terms.so gq/terms/*
CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .