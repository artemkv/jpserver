# JPServer

Play Japanese games and see an instant vocabulary list on your phone or tablet!

Starts web server with the endpoint allowing to retrieve the screenshot of the primary monitor.
Intended to work in tandem with jpterminal mobile app.

## Techinical requirements

- Windows 8 or higher
- DirectX 11.1+ compatible GPU
- Go 1.26 or higher

## How to run

- Start jpserver.exe
- Scan the QR code in the app
- Make sure your mobile phone and your laptop are connected to the same WiFi
- Play!

## Environment Variables

You can provide values for environment variables by creating '.env' (exactly like that, i.e. dot env) file next to the executable. For example:

```
JPSERVER_HOST=192.168.0.13
JPSERVER_PORT=9999
JPSERVER_CODE=824230
```

`JPSERVER_HOST`

Optional, allows to configure the interface to listen on (IP address of your network).
But default, we try to auto-detect your local network IP and use it.
If autodetect fails, defaults to empty string (same as "0.0.0.0", meaning "listen on all interfaces").
If you want to set it manually, use `ipconfig` to see the information about your network interfaces.

`JPSERVER_PORT`

Optional, allows to configure the port to listen on; defaults to "9999".
Change this if, for some reason, the port is already in use.

`JPSERVER_CODE`

Allows to specify fixed access code.
If not provided, the new access code is generated every time randomly.

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
