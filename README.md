# JPServer

**Play Japanese games on PC and see an instant vocabulary break-down on your phone or tablet!**

Starts web server with the endpoint allowing to retrieve the screenshot of the primary monitor.

Intended to work in tandem with jpterminal mobile app.

## Techinical requirements

- Windows 8 or higher
- DirectX 11.1+ compatible GPU
- Go 1.26 or higher

## How to run

- Make sure your mobile phone and your laptop are on the same WiFi
- Start `jpserver.exe` on your Windows PC
- Scan the QR code in the mobile app
- Play!

## Configuration

We try to auto-detect the most suitable configuration parameters, but you can override the defaults using the environment variables.

You can provide values for environment variables by creating `.env` (exactly like that, i.e. "dot env") file next to the executable. For example:

```
JPSERVER_HOST=192.168.0.13
JPSERVER_PORT=9999
JPSERVER_CODE_LENGTH=20
JPSERVER_CODE=824230
```

`JPSERVER_HOST`

_Optional._ Allows to configure the interface to listen on (the IP address of your network).

But default, we try to auto-detect your local network IP and use it.

If autodetect fails, defaults to empty string (same as `0.0.0.0`, meaning "listen on all interfaces").

If you want to set it manually, use `ipconfig` to see the information about your network interfaces.

`JPSERVER_PORT`

_Optional._ Allows to configure the port to listen on; defaults to the first available port between `9991` and `9999`.
Change this if, for some reason, the whole range is already in use.

`JPSERVER_CODE_LENGTH`

_Optional._ Allows to specify the length of the randomly generated access code. Default is 8 digits. If you are using `JPSERVER_CODE`, the value is ignored.

`JPSERVER_CODE`

_Optional._ Allows to specify a fixed access code. If not provided, the new random access code is generated every time. The length of the generated code is controlled by `JPSERVER_CODE_LENGTH`.

## API (for integrating your own tools)

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
