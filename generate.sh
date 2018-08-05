# we are using the gogoslick protobuf generation
protoc --gogoslick_out=plugins=grpc:. proto/space.proto proto/command.proto

# this is the default protobuf protocol
# protoc --go_out=plugins=grpc:. proto/space.proto proto/command.proto

