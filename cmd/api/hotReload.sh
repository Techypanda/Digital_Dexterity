#!/bin/sh

swag init
find /app/ | entr -r -n sh -c 'cd /app/cmd/api; go build /app/cmd/api; ./api'