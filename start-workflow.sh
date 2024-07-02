#! /usr/bin/env bash
set -e


echo "Checking server health..."
curl 'http://localhost:7999/healthz'
echo ""
echo ""

echo "Starting workflow..."
RESULT=$(curl "http://localhost:7999/workflows" \
  --silent \
  --fail-with-body \
  -X PUT \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "PostgreSQLDatabasesPut",
    "input": {
      "message": "Hello, Dapr!"
    }
  }')

echo "$RESULT"
echo ""

ID=$(jq -r '.id' <<< "$RESULT")
echo "Started workflow with id: $ID"

while true; do
  echo "Checking workflow status..."

  RESULT=$(curl "http://localhost:7999/workflows/$ID" \
    --silent \
    --fail-with-body \
    -X GET)
  echo "$RESULT"
  echo ""

## Example of status values...
#
# const (
# 	StatusRunning Status = iota
# 	StatusCompleted
# 	StatusContinuedAsNew
# 	StatusFailed
# 	StatusCanceled
# 	StatusTerminated
# 	StatusPending
# 	StatusSuspended
# 	StatusUnknown
# )

  if [[ $(jq -r '.status' <<< "$RESULT") -eq "1" ]]; then
    echo "Workflow completed!"
    break
  fi

  if [[ $(jq -r '.status' <<< "$RESULT") -eq "3" ]]; then
    echo "Workflow failed!"
    break
  fi

  if [[ $(jq -r '.status' <<< "$RESULT") -eq "4" ]]; then
    echo "Workflow canceled!"
    break
  fi

  if [[ $(jq -r '.status' <<< "$RESULT") -eq "5" ]]; then
    echo "Workflow terminated!"
    break
  fi

  sleep 3
done