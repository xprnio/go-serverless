# Hello World function
This is an example function that runs in a `node:alpine` container.

When calling this function, the following request body is expected:
```json
{
  "message": "some kind of a message"
}
```

The function will then return the message from the request as such:
```json
{
  "success": true,
  "data": {
    "message": "some kind of a message"
  }
}
```
