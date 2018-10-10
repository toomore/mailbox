#!/usr/bin/env bash
docker pull alpine:3.8
docker pull golang:1.11.1-alpine3.8
docker build -t toomore/mailbox:base ./
