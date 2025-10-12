#!/bin/bash
set -e

# Create config directory if it doesn't exist
export TALOSCONFIG_DIR="./test/configs/"
export TALOSCONFIG="./test/config"
mkdir -p $TALOSCONFIG_DIR

setup() {
    echo "Performing setup..."

    # Build talswitcher
    echo "========================================================================================="
    echo "Building talswitcher..."
    echo "========================================================================================="
    go build -o talswitcher

    # Set up talos clusters
    echo "========================================================================================="
    echo "Setting up Talos clusters..."
    echo "========================================================================================="
    talosctl cluster create --name=test-cluster-1 --talosconfig=./configs/talos1.yaml --cidr=10.5.0.0/24 &
    talosctl cluster create --name=test-cluster-2 --talosconfig=./configs/talos2.yaml --cidr=10.6.0.0/24 &
    wait

    echo "========================================================================================="
    echo "Setup completed."
    echo "========================================================================================="
}

cleanup() {
    echo "Performing cleanup..."

    # Stop and delete talos clusters
    echo "Stopping and deleting test-cluster-1..."
    talosctl cluster destroy --name=test-cluster-1 || true

    echo "Stopping and deleting test-cluster-2..."
    talosctl cluster destroy --name=test-cluster-2 || true

    # Remove config directory
    echo "Removing configs directory..."
    rm -rf ./configs/

    echo "Cleanup completed."
}

run_tests() {
    echo "========================================================================================="
    echo "Running tests..."
    echo "========================================================================================="

    echo "====> Switching to test-cluster-1..."
    ./talswitcher context test-cluster-1

    echo "====> Validating cluster switch to test-cluster-1..."
    talosctl get members -n test-cluster-1-controlplane-1

    echo "====> Switching to test-cluster-2..."
    ./talswitcher ctx test-cluster-2

    echo "====> Validating cluster switch to test-cluster-2..."
    talosctl get members -n test-cluster-2-controlplane-1

    echo "====> Attempting to list members of test-cluster-1..."
    talosctl get members -n test-cluster-1-controlplane-1 && exit 1 || echo "This was supposed to fail! We're good."

    echo "====> Switch to previous context..."
    ./talswitcher ctx -

    echo "====> Validating cluster switch to test-cluster-1..."
    talosctl get members -n test-cluster-1-controlplane-1

    echo "========================================================================================="
    echo "Tests completed successfully!"
    echo "========================================================================================="
}

# Set up trap to ensure cleanup happens on exit
trap cleanup EXIT

# Main execution
echo "Starting E2E tests..."
setup
run_tests
echo "E2E tests completed successfully!"
