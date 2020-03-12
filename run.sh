#!/bin/bash

go build -buildmode=plugin -o plugins/posts.so gq/posts/*
go build -buildmode=plugin -o plugins/translate.so gq/translate/*
go run main.go