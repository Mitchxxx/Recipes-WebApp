#!/bin/bash
while IFS= read -r thread; do
    [[ -z "$thread" || "$thread" =~ ^# ]] && continue

    printf "\n%s\n"

    curl -sS -X POST http://localhost:5050/parse \
    -H "Content-Type: application/json"\
    -d "$(jq -nc --arg url "$thread" '{url:$url}')"
done < "threads"