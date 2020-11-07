# MQTT

[![Apache-2.0 License](https://img.shields.io/github/license/htdvisser/mqtt.svg?style=flat-square)](https://github.com/htdvisser/mqtt/blob/master/LICENSE) [![GitHub stars](https://img.shields.io/github/stars/htdvisser/mqtt.svg?logo=github&style=flat-square)](https://github.com/htdvisser/mqtt/stargazers) [![GitHub forks](https://img.shields.io/github/forks/htdvisser/mqtt.svg?logo=github&style=flat-square)](https://github.com/htdvisser/mqtt/network/members) [![GitHub release](https://img.shields.io/github/release/htdvisser/mqtt.svg?logo=github&style=flat-square)](https://github.com/htdvisser/mqtt/releases) [![Go reference](https://img.shields.io/badge/go-reference-blue?style=flat-square)](https://pkg.go.dev/htdvisser.dev/mqtt)

> Package `htdvisser.dev/mqtt` implements MQTT 3.1.1 and MQTT 5.0 packet types as well as a reader and a writer.

## Background

MQTT is a publish/subscribe messaging transport protocol often used in Internet of Things communication. The MQTT specifications are maintained by [OASIS](https://www.oasis-open.org/).

## Goals and Non-Goals

The goal of this library is to provide basic MQTT packet types, as well as implementations for reading and writing those packets. This library aims to implement version [3.1.1](https://docs.oasis-open.org/mqtt/mqtt/v3.1.1/mqtt-v3.1.1.html) and version [5.0](https://docs.oasis-open.org/mqtt/mqtt/v5.0/mqtt-v5.0.html) of the specification, with limited support for version 3.1.

This library does not -- and will not -- implement a client or server (broker), but it can be used by client or server implementations.

## Install

```sh
go get -u htdvisser.dev/mqtt
```

## Usage

See the examples on [pkg.go.dev](https://pkg.go.dev/htdvisser.dev/mqtt).

## Contributing

See [CONTRIBUTING.md](.github/CONTRIBUTING.md) and [CODE_OF_CONDUCT.md](.github/CODE_OF_CONDUCT.md).

## License

[Apache-2.0](LICENSE) © Hylke Visser
