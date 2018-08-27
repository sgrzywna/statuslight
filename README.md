# statuslight

Simple web service to set Mi-Light lamp colour depending on received status. Uses [milightd](https://github.com/sgrzywna/milightd) to control the lamp.

For instance, can be used together with CI to show status of tests with appropriate colour.

## Build

To build service run:

```bash
make
```

## Start service

To connect to `milightd` running at localhost and port 8080, and to listen to commands at port 8888:

```bash
./statuslight -mihost 127.0.0.1 -miport 8080 -port 8888
```

To see all available command line switches run:

```bash
./statuslight -h
```

## Set status

API is [documented](api/swagger.yaml) with Swagger specification.

For example, to set single status:

```bash
curl -X POST "http://127.0.0.1:8888/api/v1/status" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"state\": true, \"statusId\": \"string\"}"
```
