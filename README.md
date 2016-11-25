# NPM registry

[![Build Status](https://travis-ci.org/gillesdemey/npm-registry.svg?branch=master)](https://travis-ci.org/gillesdemey/npm-registry)
[![Go Report Card](https://goreportcard.com/badge/github.com/gillesdemey/npm-registry)](https://goreportcard.com/report/github.com/gillesdemey/npm-registry)

A very simple NPM registry proxy and backup server designed to be easy, fast and stable.

## Developing

This project requires a few optional dependencies to make development easier and more enjoyable.

You can skip this step if you feel more comfortable with other tools.

### Dependencies

Install the dependencies with `make deps`

`bra` is used to watch for code changes and recompile the server.

Install it with `go get -u github.com/Unknwon/bra`

Run the server with `make run`

### Testing

Run `make test`

You can also manually test the implementations using the `npm` cli utility.

ie. `npm ping --registry http://127.0.0.1:8080/`

## Acknowledgments

Inspiration was taken from other comparable community efforts:

* [https://github.com/dickeyxxx/npm-register](https://github.com/dickeyxxx/npm-register)
* [https://github.com/rlidwka/sinopia](https://github.com/rlidwka/sinopia)

## License

MIT
