#!/usr/bin/env bash
docker pull alpine:3.6
docker pull golang:1.9.0-alpine3.6
docker build -t toomore/mailbox:base ./
