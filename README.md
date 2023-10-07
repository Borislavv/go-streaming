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
- **STREAMING_SERVER_HOST** is an WebSocket server serving host.
- **STREAMING_SERVER_PORT** is an WebSocket server serving port.
- **STREAMING_SERVER_TRANSPORT_PROTOCOL** is an WebSocket server transport protocol.
  If you are not concerned about the loss part of packets and this is not a problem for you, then use the UDP,
  because this will give you a performance gain (due to the server will not check of packages number and them ordering).
  Otherwise, if your data needs to be in safe, and you cannot afford to lose it, use the TCP.

### Database
- **MONGO_URI** is a simple MongoDb DSN string for connect to database.
- **MongoDb** is a name of database into the MongoDb.

### Application
- **JWT_SECRET_SALT** is a secret string which further will convert to slice of bytes and will be provided
  as a salt for signature the jwt tokens.
  If this variable an empty or omitted then will be generated random salt which is alive while application instance
  is running (note: you cannot get access to this value, and it will not be provided anywhere as an output. If
  you need access to this value, and you need to share this value, then sat up your own secret string).
- **UPLOADER_TYPE** is an uploading strategy which will be used for upload files on the server.
  1. '**muiltipart_form**' is a strategy which used builtin sugar approach. It will be parsing a whole file into the
            memory (if a file more than InMemoryFileSizeThreshold, it will be saved on the disk, otherwise, it will be
            loaded in the RAM).
  2. '**muiltipart_part**' is a strategy which used lower level implementation which based on the reading by parts
      from raw form data.
  If you care of application performance (speed of uploading directly) and you have enough RAM, then use
  the 'muiltipart_form' approach and increase the value of InMemoryFileSizeThreshold variable.
  Otherwise, use 'muiltipart_part' because it takes a much lower RAM per file uploading.
  For example: for upload the file which weight is 50mb. it will take around 10mb. of your RAM.

// TODO The Swagger docs. is not implemented yet :( sorry

// I'm already working on it for you! :)

At the moment, you already can surf the address: `http://0.0.0.0:8000/` in order to see the result.
