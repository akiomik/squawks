get-old-tweets
==============

[![Go](https://github.com/akiomik/get-old-tweets/actions/workflows/go.yml/badge.svg)](https://github.com/akiomik/get-old-tweets/actions/workflows/go.yml)

A cli tool to download old tweets on twitter.

## Usage

```
Usage:
  get-old-tweets --out FILENAME [flags]

Flags:
      --from string         find tweets sent from a certain user
  -h, --help                help for get-old-tweets
  -o, --out string          output csv filename (required)
  -q, --query string        query text to search
      --since string        find tweets since a certain day (e.g. 2014-07-21)
      --to string           find tweets sent in reply to a certain user
      --until string        find tweets until a certain day (e.g. 2020-09-06)
      --user-agent string   set custom user-agent
  -v, --version             version for get-old-tweets
```

## Build

```sh
make build
```

## Test

```sh
make test
```
