#!/bin/bash

go build -buildmode=plugin -o plugins/posts.so gq/posts/*
go build -buildmode=plugin -o plugins/translate.so gq/translate/*
go build -buildmode=plugin -o plugins/terms.so gq/terms/*
go build -buildmode=plugin -o plugins/menu.so gq/menu/*
go build -buildmode=plugin -o plugins/auth.so gq/auth/*
go build -buildmode=plugin -o plugins/files.so gq/files/*
go run main.go