#!/usr/bin/env sh

if [ -n "$(ls -A /usr/local/share/ca-certificates/*.crt 2> /dev/null)" ]; then
  update-ca-certificates
fi

exec "$@"