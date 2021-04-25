#!/usr/bin/env bash
docker pull alpine:3.13.5
docker pull golang:1.16.3-alpine3.13
docker build -t toomore/mailbox:base ./
