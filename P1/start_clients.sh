#!/bin/bash

# Function to run clients
run_clients() {
    local num_clients=$1
    local task_size=$2
    local delay=$3

    echo "Starting $num_clients clients with task size $task_size and delay $delay seconds..."
    for i in $(seq 1 $num_clients); do
        make run-client ARGS="$task_size" &
        CLIENT_PIDS[$i]=$!
        sleep $delay  # Add delay between client starts
    done
}

# Run 50 clients with 50-second tasks and 0.5-second delay
run_clients 50 50 0.5

# # Immediately run 100 clients with 2-second tasks and no delay
# run_clients 100 2 0

# Wait for all clients to finish
echo "Waiting for all clients to finish..."
wait
echo "All clients have finished."