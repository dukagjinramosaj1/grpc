#!/bin/bash


#AT ROOT DIRECTORY 
# export PATH="$PATH:$(go env GOPATH)/bin"


protoc greet/greetpb/greet.proto --go_out=plugins=grpc:.


protoc calculator/calculatorpb/calculator.proto --go_out=plugins=grpc:.