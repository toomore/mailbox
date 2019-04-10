#!/usr/bin/env bash
docker pull alpine:3.9
docker pull golang:1.12.3-alpine3.9
docker build -t toomore/mailbox:base ./
