#!/bin/sh

CMD="$1"
shift
if [ -z "$CMD" ]
then
  CMD="s"
fi

go run tie.prodigy9.co "$CMD"
