protoc --gofast_out=. proto/space.proto proto/command.proto
protoc -I=proto/ --gofast_out=plugins=grpc:proto/ proto/rafoe-server.proto

