# Welcome to the go-toolkit!
This repository contains many go packages, including:

- [auth](./auth)
- [httprouter](./httprouter)
- [log](./log)
- aws services
  - [config](./service/aws/config)
  - [sqs](./service/aws/sqs)
- ... and more

This monorepo was created to improve collaboration and productivity between `Platform Core`. 
By having all our code in one place, we can share ideas, find bugs and fix them more easily.

> It is not a framework but rather a set of simple utilities that 
> can be used independently of each other

## Getting started

The best way to get started working with the toolkit is to use `go get` to add the
package and desired service clients to your Go dependencies explicitly.

```shell
go github.com/pomelo-la/go-toolkit/service/aws/config
# or
go github.com/pomelo-la/go-toolkit/service/aws/sqs
# or
go get github.com/aws/aws-sdk-go-v2/log
# etc
```

## Would you like to collaborate?

You are more than welcome, please read our contribution guidelines.

- [Code Of Conduct](./code-of-conduct.md)
- [Contributing guidelines](./CONTRIBUTING.md)
