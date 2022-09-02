#!/bin/sh

if [ -n "$1" ]
then
  times=$1
else
  times=1
fi

printf "==> 'example-app' Simulator! Will make $((times*3)) good and bad requests to local 'example-app' instance!\n"

i=1
while [ $i -le $times ]; do
  printf "===========================================================================================\n"
  printf "==> #$i: Making local requests...\n"

  GOOD_REQUEST_RESP="$(curl localhost:1234/emoji\?name=100)"
  NO_MATCH_REQUEST_RESP="$(curl localhost:1234/emoji\?name=foo)"
  EMPTY_REQUEST_RESP="$(curl localhost:1234/emoji)"

  printf "==> #$i: [GOOD_REQUEST] $GOOD_REQUEST_RESP\n"
  printf "==> #$i: [NO_MATCH_REQUEST] $NO_MATCH_REQUEST_RESP\n"
  printf "==> #$i: [EMPTY_REQUEST] $EMPTY_REQUEST_RESP\n"

  ((i++))
done
printf "===========================================================================================\n"

((i--))
printf "\n==> Done! $((i*3)) local requests sent! \o/\n"
