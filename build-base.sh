#!/usr/bin/env bash
docker pull alpine:3.19.4
docker pull golang:1.22.2-alpine3.19
docker build -t toomore/mailbox:base ./
