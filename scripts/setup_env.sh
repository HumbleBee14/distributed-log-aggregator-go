#!/bin/bash

ENV_FILE=".env"
EXAMPLE_FILE=".env.example"

if [ -f "$ENV_FILE" ]; then
    echo ".env file already exists. Skipping creation."
else
    if [ -f "$EXAMPLE_FILE" ]; then
        cp "$EXAMPLE_FILE" "$ENV_FILE"
        echo ".env file created from .env.example. Please update it with your secrets."
    else
        echo "No .env.example found. Please create one first."
    fi
fi
