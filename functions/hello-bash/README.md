# Hello World function
This is an example function written in bash that runs in an `alpine` container.

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
