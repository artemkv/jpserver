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

GET /frame?left=10

Use 10% border on the left side

GET /frame?right=10

Use 10% border on the right side

GET /frame?top=20

Use 20% border on the top

GET /frame?top=20

Use 20% border on the bottom

GET /frame?resize=2

Resize the captured image by factor of 2
