# `lunch`: A command line client for lunchy

> This project was supposed to be called mclunchface to complete the
> lunchy service, but it was too tedious to write.

This is a simple command line client and library to use
the [lunchy](https://github.com/rebasar/lunchy) service from command
line. It depends on nothing but the standard Go distribution.

## Installation

Assuming that you have a Go environment setup properly:

```bash
$ go get github.com/rebasar/lunch
$ go install github.com/rebasar/lunchy
```

## Usage

Just call without any parameters, as `lunchy` to get a list of
supported places. This information is cached, by default in a file
called `~/.lunchy.cache` for now, you need to remove the file manually
if it becomes stale.

To fetch the menu for a restaurant, just write `lunchy <alias>` where
alias is one of the aliases of the place.

Enjoy your lunch!
