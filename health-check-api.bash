#!/bin/bash

# HOST=http://localhost:8080
health_check_url="$HOST/v1/health"
health_cmd="curl -w '%{http_code}\n' -s $HOST/v1/health -o /dev/null/"

retries=0
status="$(eval $health_cmd)"

while [ "$status" -ne 200 ] && [ "$retries" -lt 10 ]
do
    echo "Failed! Health check status: $status"
    echo "Sleeping..."
    status="$(eval $health_cmd)"
    let retries++
    sleep 5
done

if [ "$status" -eq 200 ] 
then
    echo "Executing: $@"
    eval "$@"
else 
    echo "Maximum retries reached... health check status: $status"
fi

