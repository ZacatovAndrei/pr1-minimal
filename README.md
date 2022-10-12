# PR1-Minimal

This is the home to the minimal acceptance criteria version for the PR
laboratory work #1

## Info

This repository contains 2 servers in 2 folders

1. `server1` - the producer/generator server, sends a small payload to `server2`
2. `server2` - the consumer server, receives the payload and processes it a bit,
   then sends back to `server1`

## Building
One can go into each folder and run `go build` to get both servers
 
The minimal acceptance criteria didn't include containerizing the servers,
hence there are no dockerfiles.

when running the servers it is important to launch the `server2` first,
as running 'server1' first will make it crash since there is nowhere to send data

there is also an included Makefile