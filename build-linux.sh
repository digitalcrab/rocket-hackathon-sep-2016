#!/usr/bin/env bash

GOOS=linux GOARCH=amd64 go build -v -o server-linux server/*
