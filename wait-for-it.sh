#!/bin/sh

set -e

host="$1"
shift
cmd="$@"

until ; do
    >&2 echo "Service is unavailable"
    sleep 1
done

>&2 echo "Service is up - executing command"
exec $cmd
