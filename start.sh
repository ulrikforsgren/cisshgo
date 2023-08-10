#!/bin/sh

[ -r /setup.sh ] && source /setup.sh

ulimit -n 100000
exec /app/cisshgo $*
