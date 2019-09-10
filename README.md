<p align="center">
  <img src="https://user-images.githubusercontent.com/132562/64657729-aa544180-d3e9-11e9-8ba5-7e5056b27c5c.png" alt="Lumberman" />
</p>

<h1 align="center">Lumberman Web Client</h1>

<p align="center">
  <strong><a href="https://github.com/webmocha/Lumberman">Lumberman</a> web user interface</strong>
</p>

## Options

| flag | default | description |
| ---- | ------- | ----------- |
| -port | 80 | Port to server web ui |
| -server_addr | localhost:9090 | Lumberman server port |

## Install and run with Go

```sh
go get github.com/webmocha/Lumberman-web-client
```

Serve on port 8080
```sh
Lumberman-web-client -port 8080
```

specify a Lumberman server address and port

```sh
Lumberman-web-client -port 8080 -server_addr otherhost:9090
```

## Run with Docker :whale:

Have the Lumberman server running
```sh
docker run -d \
  --name lumberman \
  quay.io/webmochallc/lumberman
```

Run the web client
```sh
docker run -d \
  --link lumberman \
  -p 8080:80 \
  quay.io/webmochallc/lumberman-web-client -server_addr=lumberman:9090
```

Navigate to http://localhost:8080
