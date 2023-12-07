# go-sdk

## What is go-sdk?

go-sdk provide "plumbing" primitives for creating web applications.

## What is not go-sdk?

It is not a framework but rather a set of simple utilities that 
can be used independently of each other.

## How does it work?
The project is split internally into different modules to maintain 
each module with its dependencies; this allows different releases 
to be made for each module separately, the client downloads 
fewer transitive dependencies, and development can be done 
with better quality and speed.

The best way to get started working with the toolkit is 
to use `go get` to add the package and desired service clients 
to your Go dependencies explicitly.


    go get github.com/pomelo-la/go-toolkit/service/aws/config
    or
    go get github.com/pomelo-la/go-toolkit/logger

Below is a list of the modules provided:

- [Packages](./packages.md)
