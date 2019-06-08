#!/bin/bash

cd $(cd $(dirname $0) && pwd)

go build -o sleep

./sleep
NORMAL_EXIT_CODE=$?

echo "NORMAL_EXIT_CODE:${NORMAL_EXIT_CODE}"
