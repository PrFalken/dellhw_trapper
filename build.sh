#!/bin/bash

dir=$( dirname $0 )


go get ./...
mkdir ${dir}/dist
GOOS=linux go build -o ${dir}/dist/hardware_exporter
