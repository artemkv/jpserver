# JPServer

The web server that exposes the endpoint allowing to retrieve the screenshot of primary monitor

## Environment Variables

```
JPSERVER_HOST=192.168.0.13
JPSERVER_PORT=9999
JPSERVER_ALLOW_ORIGIN=https://localhost
JPSERVER_CODE=824230
```

`JPSERVER_HOST`

Optional, allows to configure the interface to listen on
Defaults to auto-detected local IP
If fails to autodetect, defaults to empty string (same as "0.0.0.0")

`JPSERVER_PORT`

Optional, allows to configure the port to listen on; defaults to "9999"

`JPSERVER_ALLOW_ORIGIN`

Optional, allows to configure the allowed origin for CORS; defaults to "https://localhost"
This should match the origin of the running app (reported in the log for every request)

`JPSERVER_CODE`

Allows to specify fixed access code
If not provided, the new access code is generated every time randomly

## API

`GET /frame`

Returns the screenshot (in PNG format)

`GET /frame?left=10`

Use 10% border on the left side

`GET /frame?right=10`

Use 10% border on the right side

`GET /frame?top=20`

Use 20% border on the top

`GET /frame?top=20`

Use 20% border on the bottom

`GET /frame?resize=2`

Resize the captured image by factor of 2
