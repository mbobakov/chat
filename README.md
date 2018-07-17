## Simple grpc-stream based chat
### Usage
`Terminal#1 (SERVER)`
```
▶ go run cmd/server/main.go
2018/07/17 20:25:43 main.go:23: [INFO] Launching Application with: {GRPCListen: DebugListen: Verbose:false}
2018/07/17 20:25:43 server.go:57: [INFO] Start listening on ':8080'

```
`Terminal#1 (Client alice)`
```
▶ go run cmd/client/main.go --username alice
Connecting to 127.0.0.1:8080
Connected
Happy Chatting. Your Username: 'alice'
```
`Terminal#1 (Client bob)`
```
▶ go run cmd/client/main.go --username bob
Connecting to 127.0.0.1:8080
Connected
Happy Chatting. Your Username: 'bob'
```
Now you can send messages to the chat from alice or bob