#!/bin/bash

go build -buildmode=plugin -o plugins/posts.so gq/posts/*
go build -buildmode=plugin -o plugins/translate.so gq/translate/*
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .