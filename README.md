# Streaming service

---

## Overview
![image](https://github.com/Borislavv/go-streaming/assets/50691459/7cf3b525-df6b-49e4-8339-b4f51074b55d)

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

---

## Configuration

### Api
- **API_VERSION_PREFIX** is a value which will be used as your RestAPI controllers version prefix.
  For example: {{schema}}://{{host}}:{{port}}{{ApiVersionPrefix}}/{{additionalControllerPath}}.
  Default: `/api/v1`
- **RENDER_VERSION_PREFIX** is a value which will be used as your native rendering controllers version prefix.
  For example: {{schema}}://{{host}}:{{port}}{{RenderVersionPrefix}}/{{additionalControllerPath}}.
  By default, it's an empty string.
- **STATIC_VERSION_PREFIX** is a value which will be used as your static files controllers version prefix.
  For example: {{schema}}://{{host}}:{{port}}{{StaticVersionPrefix}}/{{additionalControllerPath}}.
  By default, it's an empty string.

### Server
1. #### HTTP
   - **RESOURCES_SERVER_HOST** is an HTTP server serving host. Default: `0.0.0.0`.
   - **RESOURCES_SERVER_PORT** is an HTTP server serving port. Default: `8000`.
   - **RESOURCES_SERVER_TRANSPORT_PROTOCOL** is an HTTP server transport protocol. Default: `tcp`.
     If you are not concerned about the loss part of packets and this is not a problem for you, then use the UDP,
     because this will give you a performance gain (due to the server will not check of packages number and them ord.).
     Otherwise, if your data needs to be in safe, and you cannot afford to lose it, use the TCP.
2. #### WebSocket
   - **STREAMING_SERVER_HOST** is an WebSocket server serving host. Default: `0.0.0.0`.
   - **STREAMING_SERVER_PORT** is an WebSocket server serving port. Default: `9988`.
   - **STREAMING_SERVER_TRANSPORT_PROTOCOL** is an WebSocket server transport protocol. Default: `tcp`.
     If you are not concerned about the loss part of packets and this is not a problem for you, then use the UDP,
     because this will give you a performance gain (due to the server will not check of packages number and them ordering).
     Otherwise, if your data needs to be in safe, and you cannot afford to lose it, use the TCP.

### Database
- **MONGO_URI** is a simple MongoDb DSN string for connect to database. Default: `mongodb://mongodb:27017/streaming`.
- **MongoDb** is a name of database into the MongoDb. Default: `streaming`.

### Application
- **JWT_SECRET_SALT** is a secret string which further will convert to slice of bytes and will be provided
  as a salt for signature the jwt tokens.
  If this variable an empty or omitted then will be generated random salt which is alive while application instance
  is running (note: you cannot get access to this value, and it will not be provided anywhere as an output. If
  you need access to this value, and you need to share this value, then sat up your own secret string).
- **JWT_TOKEN_ISSUER** is an issuer of JWT token. This variable helps determine which service issued the token
  (commonly used for verify that token was created by service of your system, for example,
  if you have more than one service which able for issue a token).
- **JWT_TOKEN_ACCEPTED_ISSUERS** is a string with another JwtTokenIssuer values separated by delimiter. This values will
  be accepted while token payload verification.
- **JWT_TOKEN_ENCRYPT_ALGO** is a value which will be used as encrypt algo for encode the token. Default: `HS256`.
- **UPLOADER_TYPE** is an uploading strategy which will be used for upload files on the server. Default: `muiltipart_part`.
  1. '**muiltipart_form**' is a strategy which used builtin sugar approach. It will be parsing a whole file into the
            memory (if a file more than InMemoryFileSizeThreshold, it will be saved on the disk, otherwise, it will be
            loaded in the RAM).
  2. '**muiltipart_part**' is a strategy which used lower level implementation which based on the reading by parts
      from raw form data.
  If you care of application performance (speed of uploading directly) and you have enough RAM, then use
  the 'muiltipart_form' approach and increase the value of InMemoryFileSizeThreshold variable.
  Otherwise, use 'muiltipart_part' because it takes a much lower RAM per file uploading.
  For example: for upload the file which weight is 50mb. it will take around 10mb. of your RAM.
- **RESOURCE_FORM_FILENAME** is a value which will be used for extract a file from the form by given string. Default: `resource`.
  *Used only with the 'muiltipart_form' strategy because the 'muiltipart_part' will search the first form file.
  Be careful and don't send more than one file per request in one form.
- **MAX_UPLOADING_FILESIZE** is a threshold value which means the max. weight of uploading file in bytes. Default: `5368709120`.
  By default, it's 5gb per file.
- **IN_MEMORY_FILE_SIZE_THRESHOLD** is a threshold value which means the max. weight of uploading file in bytes
  which may be loaded in the RAM. Default: `104857600`. If file weight is more this value, than it will be loaded on the disk (slow op.).
  By default, it's 100mb per file.
- **ADMIN_CONTACT_EMAIL_ADDRESS** is a target administrator contact email address for takes a users errors reports.

### Logger
- **LOGGER_ERRORS_BUFFER_CAPACITY** is errors channel capacity. Default: `10`.
  Logger is basing on the go channels, this value will be sat up as capacity.
- **LOGGER_REQUESTS_BUFFER_CAPACITY** is requests channel capacity. Default: `10`.
  Use only when you are logging input requests/responses.

### File reader
- **FILE_READER_CHUNK_SIZE** is a value which means the size of one chunk while reading the file when streaming a resource.
  By default, it's 1mb. Default: `1048576`.

---

## Launching

At the moment, you already can surf the address: `http://0.0.0.0:8000/` in order to see the result.


---

## Not implemented:

1. The Swagger docs.
2. Tests.
3. Benchmarks.

### I'm already working on it for you! :)

