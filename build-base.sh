#!/usr/bin/env bash
docker pull alpine:3.7
docker pull golang:1.9.4-alpine3.7
docker build -t toomore/mailbox:base ./
