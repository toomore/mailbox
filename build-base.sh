#!/usr/bin/env bash
docker pull alpine:3.18.2
docker pull golang:1.20.5-alpine3.18
docker build -t toomore/mailbox:base ./
