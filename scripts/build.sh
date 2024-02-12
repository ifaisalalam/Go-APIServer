#!/usr/bin/env sh

test -f $CURR_DIR/cmd/api-server/main.go

if [ $? -eq 1 ]; then
  make generate
fi
