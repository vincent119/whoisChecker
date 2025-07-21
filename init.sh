#!/bin/bash


go mod init whoisChecker
go mod tidy
go install github.com/spf13/cobra-cli@latest
# cobra-cli init 
# cobra-cli add whois