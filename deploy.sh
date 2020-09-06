#! /usr/bin/env bash

set -ue -o pipefail
trap "echo >&2 script '${BASH_SOURCE[0]}' failed" ERR

SCRIPT=`readlink -e "${BASH_SOURCE[0]}"`
SCRIPT_DIR=`dirname "$SCRIPT"`
cd "$SCRIPT_DIR"
unset SCRIPT SCRIPT_DIR

IMAGE=page-change-mailer
CONTAINER=page-change-mailer
VOLUME=page-change-mailer

docker build --pull --force-rm --tag "$IMAGE" .

docker stop "$CONTAINER" || true
docker rm --volumes "$CONTAINER" || true

docker run \
	--env-file env \
	--mount type=volume,source="$VOLUME",destination=/home/user/data \
	--name "$CONTAINER" \
	--read-only \
	--restart always \
	"$IMAGE"
