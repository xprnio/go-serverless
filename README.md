# go-serverless - Run Docker containers as serverless funcitons

**This application is in active development**

This means that some features might not be quite polished yet.
Expect bugs, but feel free to report them here.

---

## What is this
`go-serverless` is an application that makes it easy to run Docker containers as if they were serverless functions through a simple HTTP API. It exposes a set of API endpoints to both manage these functions as well as call them.

### How do I create these serverless function images?
You start with putting together a regular Docker image. The image itself can use any language and runtime you can think of, from Go to Rust to Node.js to COBOL if you're feeling like a [Chad](https://github.com/ThePrimeagen/CHADstack). As long as you can run it in Docker, you should be good to go.

Once you have your Docker image, you should get to actually building out your serverless function. Since your serverless function container will be working over an HTTP api, it will be given the JSON request body and you are expected to return a JSON response.

When running your code, the following environment variables will be set:
```env
# REQUEST_ID will be a UUID unique to this function call
REQUEST_ID="00000000"

# CONTEXT_PATH will be based off of the REQUEST_ID
# and is what you will need to use
CONTEXT_PATH="/tmp/context/00000000"
```

Additionally, two files are also passed to the container:
```sh
# The request body passed to your function
$CONTEXT_PATH/request.json

# The resulting response of your function
$CONTEXT_PATH/response.json
```

Your container should read the `request.json` file as necessary,
[do what it do](https://youtu.be/watch?v=W8LJWis1JoI), and write the response into `response.json`. Once your container exits, the response will be passed to the function caller through the HTTP API. 

You can have a look at the [hello-world](functions/hello-world) function for inspiration.

## Getting started
There's a few ways to run `go-serverless`:
- Running it directly on your machine
- Running it directly on docker
- Running it with docker compose

### 1. Running it directly on your machine
This method assumes that you at the very least have `git` and `go` installed.

Start by cloning the directory
```sh
git clone https://github.com/xprnio/go-serverless
cd go-serverless
```
  
If you have `make` installed, you can just run `make run` which will both build and run the app
```sh
make run
```
  
If you don't have `make` installed and don't want to, you can also just build it directly
```sh
mkdir bin
go build \
  -o bin/go-serverless \
  cmd/main.go

# Start the application
./bin/go-serverless
```

### 2. Running it directly on docker
The below environment variables are optional, as they are also included in the Dockerfile and are already set to the values shown below. However, if you want to use a different volume for the request contexts, you will need to modify both the environment variables as well as the volume mount.
```sh
docker run -d \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e CONTEXT_NAME=serverless_context \
  -e CONTEXT_PATH=/tmp/context \
  -v serverless_context:/tmp/context \
  -v data:/data \
  -p 9999:9999 \
  ghcr.io/xprnio/go-serverless
```

### 3. Running it with docker compose
This method is probably the most convenient way of getting started since it makes configuring the app much more easier.

Start by grabbing the [docker-compose.yml](docker-compose.yml)
```sh
curl -o docker-compose.yml \
  https://raw.githubusercontent.com/xprnio/go-serverless/main/docker-compose.yml
```

Once that is downloaded, just run the following
```sh
docker compose up -d
```

## Using the API

### Listing functions
```sh
curl http://localhost:9999/v1/functions
```

### Listing routes
```sh
curl http://localhost:9999/v1/routes
```

### Creating functions
To create a function, you will need an image.
Currently only public registries are supported.

```sh
curl http://localhost:9999/v1/functions -d '{
  "image": "ghcr.io/xprnio/serverless-hello-world:latest"
}'
```

### Creating routes
Once you have created the function, you can route it to a specific path through it's `id`.

```sh
curl http://localhost:9999/v1/routes -d '{
  "path": "/hello-world",
  "function_id": "00000000-0000-0000-0000-000000000000"
}'
```

### Running functions
Once you have routed your function, you can call it from `/r/{path}`

```sh
curl http://localhost:9999/r/hello-world -d '{
  "message": "hello-world"
}'
```
