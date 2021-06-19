#!/bin/bash

# default arguments
BOARD=
START_PAGE=
END_PAGE=
OUTPUT_DIR=


function usage {
  echo "Usage: $0 [-b <board name>] [-s <start page>] [-e <end page>] [-o <output dir>]"
  exit 1
}

while getopts "b:s:e:o:" opt; do
  case "$opt" in
    b)
      BOARD=$OPTARG
      ;;
    s)
      START_PAGE=$OPTARG
      ;;
    e)
      END_PAGE=$OPTARG
      ;;
    o)
      OUTPUT_DIR="$OPTARG"
      ;;
    *)
      usage
      ;;
   esac
done

if [[ -z $BOARD || -z $START_PAGE || -z $END_PAGE || -z "$OUTPUT_DIR" ]]; then
  usage
fi

CON_OUTPUT_DIR=/opt/crawler/data
docker run --mount type=bind,src="$OUTPUT_DIR",target="$CON_OUTPUT_DIR" ptt-crawler \
    -board $BOARD -directory "$CON_OUTPUT_DIR" -start $START_PAGE -end $END_PAGE
