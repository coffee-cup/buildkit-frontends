# Buildkit Frontends

Repo that shows an example of how to create custom Buildkit frontends in both Go and Rust (not done yet).

## Setting up Buildkit

To use the custom frontends, you need to have Buildkit setup and running locally.

On macOS, you can use the following command to install Buildkit:

```bash
brew install buildkit
```

Mac doesn't support running buildkitd as a service, but we can create a container and set that as the daemon.

```bash
docker run --rm --privileged -d --name buildkit moby/buildkit

# Set the buildkit host to the container
export BUILDKIT_HOST=docker-container://buildkit
```

## Building the frontends

Buildkit frontends run as a GRPC server that that Buildkit can talk to. For both Go and Rust we create an image that runs our server.

There are custom Dockerfiles for each.

### Go

Build the frontend image.

```bash
docker build -t go-frontend:latest -f go-frontend.Dockerfile .

# I gave up on figuring out how to get the buildkitd container to use the local registry,
# so I just pushed to ghcr.io. You will have to customize this for your own account.
# docker tag go-frontend:latest ghcr.io/coffee-cup/buildkit-frontends:go-frontend
# docker push ghcr.io/coffee-cup/buildkit-frontends:go-frontend
```

We can then tell Buildkit to build our example using our custom gateway
frontend. We pipe the result into `docker load` to save it as an image locally.

```bash
buildctl build \
    --frontend=gateway.v0 \
    --local context=examples/node-npm-latest \
    --opt source=ghcr.io/coffee-cup/buildkit-frontends:go-frontend \
    --output type=docker,name=go-node | docker load
```

Run the image

```bash
docker run --rm -it go-node
```
