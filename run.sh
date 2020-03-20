#!/bin/bash

go build -buildmode=plugin -o plugins/posts.so gq/posts/*
go build -buildmode=plugin -o plugins/translate.so gq/translate/*
go build -buildmode=plugin -o plugins/terms.so gq/terms/*
go run main.go