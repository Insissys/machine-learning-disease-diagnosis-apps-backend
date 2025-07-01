#!/bin/sh

# Replace env vars in config template
envsubst < ./internal/config/config.yml.template > ./internal/config/config.yml

# Run the Go app
exec ./server
