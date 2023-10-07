# Video/audio streaming service

## Installation

It's a pretty simple process, first of all you need to build the application:
- You can use makefile for simplicity this process, just type: `make up` and enjoy or `make build` and check the `app/cmd/built`, there must be already `streaming` binary file, ok, let's do it again, type `make build` and go the same dir., you should to see the `legacy_streaming` binary file. This feature is implemented by `./.ops/build/build_and_rotate.sh` script, which target is build and rotate assemblies as you colud notice.
  
  ```
    make build  // will build and rotate your assemblies
    make up     // this way is suitable for local development, without compilation of binary file
  ```
  
- Or you can use the good old way with `docker-compose build` and `docker-compose up` if you wish.

  ```
    docker-compose build  // will build the appliocation with assemblies rotation
    docker-compose up     // will run the application with assemblies roration
  
    docker-compose -f docker-compose.local.yaml build // will build the appliocation for local development
    docker-compose -f docker-compose.local.yaml up    // will run the appliocation for local development
  ```
  
- Of course, don’t forget to set the environment variables according to your needs. You won’t have to look at them in the `docker-compose` file.

## Configuration

### Api
- **API_VERSION_PREFIX** is a value which will be used as your RestAPI controllers version prefix.
   For example: {{schema}}://{{host}}:{{port}}{{ApiVersionPrefix}}/{{additionalControllerPath}}
- **RENDER_VERSION_PREFIX** is a value which will be used as your native rendering controllers version prefix.
  For example: {{schema}}://{{host}}:{{port}}{{RenderVersionPrefix}}/{{additionalControllerPath}}
  By default it's an empty string.
- **STATIC_VERSION_PREFIX** is a value which will be used as your static files controllers version prefix.
  For example: {{schema}}://{{host}}:{{port}}{{StaticVersionPrefix}}/{{additionalControllerPath}}
  By default it's an empty string.

### Server
#### HTTP
- **RESOURCES_SERVER_HOST** is an HTTP server serving host.
- **RESOURCES_SERVER_PORT** is an HTTP server serving port.
- **RESOURCES_SERVER_TRANSPORT_PROTOCOL** is an HTTP server transport protocol.
  If you are not concerned about the loss part of packets and this is not a problem for you, then use the UDP,
  because this will give you a performance gain (due to the server will not check of packages number and them ord.).
  Otherwise, if your data needs to be in safe, and you cannot afford to lose it, use the TCP.
#### WebSocket

// TODO The Swagger docs. is not implemented yet :( sorry

// I'm already working on it for you! :)

At the moment, you already can surf the address: `http://0.0.0.0:8000/` in order to see the result.
