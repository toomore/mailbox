#!/usr/bin/env bash
docker pull alpine:3.7
docker pull golang:1.10.0-alpine3.7
docker build -t toomore/mailbox:base ./
