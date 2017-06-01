#!/usr/bin/env bash
docker pull alpine:latest
docker pull golang:alpine
docker build -t toomore/mailbox:base ./
