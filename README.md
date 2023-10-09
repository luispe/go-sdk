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

Please read [CODE OF CONDUCT](./code-of-conduct.md)

> NOTE
> 
> For privacy and security reasons the Pomelo repositories are private, 
> please make the following settings to avoid errors during the import of the different packages

First step create a [gitHub token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens) 
and allow reading and writing of repositories.

In your preferred terminal run

```shell
touch ~/.netrc
echo "machine github.com login <YOUR_GITHUB_USER> password <THE_PREVIOUSLY_GENERATED_TOKEN> >> ~/.netrc
```

One final thing that you’ll have to do is set a `GOPRIVATE` environment variable. 
This contains a comma-separated list of module prefixes. Save this value to your `~/.bashrc` or `~/.zshrc`.

e.g
```shell
 echo "export GOPRIVATE="github.com/pomelo-la/,github.com/pomelo-la/*"" >> "$HOME"/.zshrc
```

---

Now that you have everything ready

The best way to get started working with the toolkit is to use `go get` to add the
package and desired service clients to your Go dependencies explicitly.

```shell
go github.com/pomelo-la/go-toolkit/service/aws/config
# or
go github.com/pomelo-la/go-toolkit/service/aws/sqs
# or
go github.com/pomelo-la/go-toolki/log
# etc
```

## Getting Help

* [GitHub discussions](https://github.com/pomelo-la/go-toolkit/discussions) - For ideas, RFCs & general questions
* [GitHub issues](https://github.com/pomelo-la/go-toolkit/issues/new/choose) – For bug reports & feature requests
* [Usage examples](https://github.com/pomelo-la/go-toolkit-examples)

### Contributing

If you are interested in contributing to the SDK, please take a look at [CONTRIBUTING](./CONTRIBUTING.md)