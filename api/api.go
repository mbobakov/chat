package api

//go:generate sh -c "cd .. && protoc -I api api/chat.proto  --go_out=plugins=grpc:api"

const userContextKey = "chatusername"
