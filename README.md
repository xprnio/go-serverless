# go-serverless - Run Docker containers as serverless funcitons

**This application is in active development**
This means that some features might not be quite polished yet.
Expect bugs, but feel free to report them here.

### Getting started
```sh
git clone https://github.com/xprnio/go-serverless
cd go-serverless

make build
```

### Listing functions
```sh
curl http://localhost:9999/v1/functions
```

### Listing routes
```sh
curl http://localhost:9999/v1/routes
```

### Creating functions
```sh
curl http://localhost:9999/v1/functions -sd '{
  "image": "ghcr.io/xprnio/serverless-hello-world:latest"
}'
```

### Creating routes
```sh
curl http://localhost:9999/v1/routes -sd '{
  "function_id": "...",
  "path": "/hello-world"
}'
```

### Running functions
```sh
curl http://localhost:9999/r/hello-world -sd '{
  "message": "hello-world"
}'
```
