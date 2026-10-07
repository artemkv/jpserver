# JPServer

The web server that exposes the endpoint allowing to retrieve the screenshot of primary monitor

## Environment Variables

```
JPSERVER_PORT=:9999
JPSERVER_ALLOW_ORIGIN=https://localhost
```

## API

GET /frame

Returns the screenshot (in PNG format)
