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
    sudo -E talosctl cluster create dev --name=talos-cluster-1 --talosconfig=$TALOSCONFIG_DIR/talos1.yaml --cidr=10.5.0.0/24 --wait=false &
    sudo -E talosctl cluster create dev --name=talos-cluster-2 --talosconfig=$TALOSCONFIG_DIR/talos2.yaml --cidr=10.6.0.0/24 --wait=false &
    wait

    echo "========================================================================================="
    echo "Setup completed."
    echo "========================================================================================="
}

cleanup() {
    echo "Performing cleanup..."

    # Stop and delete talos clusters
    echo "Stopping and deleting talos-cluster-1..."
    talosctl cluster destroy --name=talos-cluster-1 || true

    echo "Stopping and deleting talos-cluster-2..."
    talosctl cluster destroy --name=talos-cluster-2 || true

    # Remove config directory
    echo "Removing configs directory..."
    rm -rf ${TALOSCONFIG_DIR:?}/

    echo "Cleanup completed."
}

run_tests() {
    echo "========================================================================================="
    echo "Running tests..."
    echo "========================================================================================="

    echo "====> Switching to talos-cluster-1..."
    ./talswitcher context talos-cluster-1

    echo "====> Validating cluster switch to talos-cluster-1..."
    talosctl get members -n talos-cluster-1-controlplane-1

    echo "====> Switching to talos-cluster-2..."
    ./talswitcher ctx talos-cluster-2

    echo "====> Validating cluster switch to talos-cluster-2..."
    talosctl get members -n talos-cluster-2-controlplane-1

    echo "====> Attempting to list members of talos-cluster-1..."
    talosctl get members -n talos-cluster-1-controlplane-1 && exit 1 || echo "This was supposed to fail! We're good."

    echo "====> Switch to previous context..."
    ./talswitcher ctx -

    echo "====> Validating cluster switch to talos-cluster-1..."
    talosctl get members -n talos-cluster-1-controlplane-1

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
