#!/bin/bash

set -e

# Validate arguments
if [ $# -lt 2 ]; then
    echo "Usage: $0 <service> <config_file>"
    echo "  service: api, platform, or landing"
    echo "  config_file: path to configuration file"
    exit 1
fi

SERVICE=$1
CONFIG_FILE=$2

# Validate service type
if [ "$SERVICE" != "api" ] && [ "$SERVICE" != "platform" ] && [ "$SERVICE" != "landing" ] && [ "$SERVICE" != "hostname" ]; then
    echo "Error: service must be one of: api, platform, landing, hostname"
    exit 1
fi

# Validate config file exists
if [ ! -f "$CONFIG_FILE" ]; then
    echo "Error: config file not found: $CONFIG_FILE"
    exit 1
fi

echo "Starting Zewi $SERVICE with config: $CONFIG_FILE"

# Show version
./zewi version

# Run migrations only for services that need database
if [ "$SERVICE" = "api" ]; then
    echo "Running database migrations..."
    ./zewi migrate up -c "$CONFIG_FILE"
fi

# Start the service
echo "Starting $SERVICE service..."
exec ./zewi "$SERVICE" -c "$CONFIG_FILE"
