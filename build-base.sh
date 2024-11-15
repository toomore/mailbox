#!/usr/bin/env bash
docker pull alpine:3.20.2
docker pull golang:1.22.2-alpine3.20
docker build -t toomore/mailbox:base ./
