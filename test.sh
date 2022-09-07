#!/bin/sh


SMOKE=`go env GOPATH`/bin/smoke

if [ ! -x "$SMOKE" ]
then
  go install -v github.com/chakrit/smoke@latest
fi

cd tests && "$SMOKE" tests.yml "$@"
