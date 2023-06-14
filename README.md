# Server-Assignment-TikTokTechImmersion
Server assignment for  TikTok TechImmersion 2023. Email: xiet0005@e.ntu.edu.sg

This project is one of Mr. Xie Tianze's project assignments for the TikTokTechImmersion 2023.

This project is an implementation of a simple gRPC based chat service. This service uses a SQLite3 database to store and retrieve chat messages. The main features of this service include sending messages to a chat and pulling messages from a chat.

This is the first time I have worked independently on a project, no matter what programming language it is. It was a very rewarding and painful start. There are a few things I haven't implemented yet compared to the official demo and I will continue to learn.

## Getting Started

Follow the instructions below to get this project running on your local machine for development and testing purposes.

### Prerequisites

- Go v1.15 or above
- SQLite3
- Google Protocol Buffers

### Installing

- Clone the repository
```
git clone https://github.com/yourusername/xtz-Server-Assignment-TikTokTechImmersion.git
```

- Navigate to the project directory
```
cd xtz-Server-Assignment-TikTokTechImmersion
```

- Install the necessary Go packages
```
go get -u google.golang.org/grpc
go get -u google.golang.org/protobuf/cmd/protoc-gen-go
go get -u github.com/mattn/go-sqlite3
go get -u github.com/stretchr/testify/require
```

- Compile the protobuf files
```
protoc -I api/ --go_out=plugins=grpc:api api/api.proto
```

### Running

- To start the server, run the following command
```
go run server/server.go
```

The server will start and listen on port 8888.

## Testing

To run the tests, navigate to the `test` directory and run the following command
```
go test -v
```

This will run all the unit and benchmark tests and output the results to the console.

## Structure

The project structure is as follows:

```
root/
├── api/
│   ├── api.proto
├── client/
│   ├── client.go
├── server/
│   ├── server.go
├── test/
│   ├── rpc-server_test.go
│   ├── rpc-server_test_benchmark.go
```

- `api/`: This directory contains the Protobuf definitions for the gRPC service.
- `client/`: This directory contains a simple command-line client to interact with the gRPC service.
- `server/`: This directory contains the gRPC server implementation.
- `test/`: This directory contains the unit and benchmark tests for the gRPC server.

## Client Usage

Start the client by running the following command
```
go run client/client.go
```

### Send Message

To send a message, use the following format
```
send <chat> <message> <sender>
```
Example:
```
send chat1 hello user1
```

### Pull Messages

To pull messages, use the following format
```
pull <chat> <cursor> <limit>
```
Example:
```
pull chat1 0 10
```

## Benchmark Testing

Includes a benchmark test to simulate multiple concurrent users sending and extracting messages from the service, currently preset at 20 concurrent threads.

You can run this test by navigating to the `test` directory and running the following command
```
go test -bench=.
```

## Built With

- [Go](https://golang.org/) - The programming language used
- [gRPC](https://grpc.io/) - The RPC framework used
- [SQLite3](https://www.sqlite.org/index.html) - The database system used
- [Google Protocol Buffers](https://developers.google.com/protocol-buffers) - The language-neutral, platform-neutral, extensible mechanism for serializing structured data used
- [Testify](https://github.com/stretchr/testify) - The Go testing toolkit used

## Acknowledgments

This project is one of Mr. Xie Tianze's project assignments for the TikTokTechImmersion 2023.

