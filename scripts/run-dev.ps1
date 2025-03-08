# Stop and remove existing container if it exists
podman container stop crate-data-dev 2>$null
podman container rm crate-data-dev 2>$null

# Run development container with host network
podman run -it `
    --name crate-data-dev `
    --network host `
    -v ${PWD}:/app `
    crate-data-env:with-deps