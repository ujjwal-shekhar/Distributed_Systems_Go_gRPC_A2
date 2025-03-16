#!/bin/bash

# Check if a policy is provided as an argument
if [ -z "$1" ]; then
    echo "Usage: ./start_servers.sh <policy>"
    echo "Available policies: round_robin, least_loaded, pick_first"
    exit 1
fi

POLICY=$1

# Validate the policy
if [[ "$POLICY" != "round_robin" && "$POLICY" != "least_loaded" && "$POLICY" != "pick_first" ]]; then
    echo "Invalid policy: $POLICY"
    echo "Available policies: round_robin, least_loaded, pick_first"
    exit 1
fi

# Start the Load Balancer
echo "Starting Load Balancer with policy: $POLICY..."
make run-lb POLICY=$POLICY &
LB_PID=$!
sleep 2

# Start 10 Servers
echo "Starting Servers..."
for i in {1..10}; do
    make run-server &
    SERVER_PIDS[$i]=$!
    sleep 1
done

# Wait for servers to initialize
sleep 5

echo "Load Balancer and Servers are running."
echo "Load Balancer PID: $LB_PID"
echo "Server PIDs: ${SERVER_PIDS[*]}"