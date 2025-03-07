#!/bin/bash
set -e

# Create config directory if it doesn't exist
export TALOSCONFIG_DIR="./configs/"
export TALOSCONFIG="./config"
export TALSWITCHER_LOG_LEVEL="debug"
mkdir -p $TALOSCONFIG_DIR

cleanup() {
    echo "Performing cleanup..."

    # Stop and delete minikube clusters
    echo "Stopping and deleting test-cluster-1..."
    talosctl cluster destroy --name=test-cluster-1 || true

    echo "Stopping and deleting test-cluster-2..."
    talosctl cluster destroy --name=test-cluster-2 || true

    # Remove config directory
    echo "Removing configs directory..."
    rm -rf ./configs/

    echo "Cleanup completed."
}

setup() {
    echo "Performing setup..."

    # Build talswitcher
    echo "Building talswitcher..."
    go build -o talswitcher

    # Set up Kubernetes clusters
    talosctl cluster create --name=test-cluster-1 --talosconfig=./configs/talos1.yaml --cidr=10.5.0.0/24 &
    talosctl cluster create --name=test-cluster-2 --talosconfig=./configs/talos2.yaml --cidr=10.6.0.0/24 &
    wait

    echo "Setup completed."
}

run_tests() {
    echo "Running tests..."

    # Test cluster switching
    echo "Switching context to test-cluster-1..."
    ./talswitcher test-cluster-1

    echo "Validating cluster switch to test-cluster-1..."
    talosctl get members -n test-cluster-1-controlplane-1

    echo "Validating cluster switch to test-cluster-2..."
    ./talswitcher test-cluster-2

    echo "Validating cluster switch to test-cluster-2..."
    talosctl get members -n test-cluster-1-controlplane-2
    
    echo "Attempting to list members of test-cluster-1..."
    talosctl get members -n test-cluster-1-controlplane-1 && exit 1 || echo "This was supposed to fail! We're good."

    echo "Tests completed successfully!"
}

# Set up trap to ensure cleanup happens on exit
trap cleanup EXIT

# Main execution
echo "Starting E2E tests..."
setup
run_tests
echo "E2E tests completed successfully!"
