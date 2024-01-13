#!/bin/bash

function read_request {
  cat "$CONTEXT_PATH/request.json"
}

function write_response {
  echo "$1" > "$CONTEXT_PATH/response.json"
}

message=`read_request | jq -r '.message'`

if [[ "$message" == "" ]]; then
  write_response "{
    \"success\": false,
    \"message\": \"missing message\"
  }"
  exit 1
fi

# Send response
write_response "{
  \"success\": true,
  \"message\": \"$message\"
}"
