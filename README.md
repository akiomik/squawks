get-old-tweets
==============

[![Go](https://github.com/akiomik/get-old-tweets/actions/workflows/go.yml/badge.svg)](https://github.com/akiomik/get-old-tweets/actions/workflows/go.yml)
[![CodeQL](https://github.com/akiomik/get-old-tweets/actions/workflows/codeql.yml/badge.svg)](https://github.com/akiomik/get-old-tweets/actions/workflows/codeql.yml)

A cli tool to get old tweets on twitter (inspired by [Jefferson-Henrique/GetOldTweets-python](https://github.com/Jefferson-Henrique/GetOldTweets-python)).

## Install

- Download the latest release from [the release page](https://github.com/akiomik/get-old-tweets/releases/latest). Supported platforms are Windows, macOS and linux.
- Extract the downloaded `.tar.gz` file.

## Usage

```
Usage:
  get-old-tweets --out FILENAME [flags]

Flags:
      --exclude strings     exclude tweets by type of tweet [hashtags|nativeretweets|retweets|replies] (default [])
      --filter strings      find tweets by type of account or tweet [verified|follows|media|images|twimg|videos|periscope|vine|consumer_video|pro_video|native_video|links|hashtags|nativeretweets|retweets|replies|safe|news] (default [])
      --from string         find tweets sent from a certain user
      --geocode string      find tweets sent from certain coordinates (e.g. 35.6851508,139.7526768,0.1km)
  -h, --help                help for get-old-tweets
      --include strings     include tweets by type of tweet [hashtags|nativeretweets|retweets|replies] (default [])
      --lang string         find tweets by a certain language (e.g. en, es, fr)
      --near string         find tweets nearby a certain location (e.g. tokyo)
  -o, --out string          output csv filename (required)
  -q, --query string        query text to search
      --since string        find tweets since a certain day (e.g. 2014-07-21)
      --to string           find tweets sent in reply to a certain user
      --top                 find top tweets
      --until string        find tweets until a certain day (e.g. 2020-09-06)
      --url string          find tweets containing a certain url (e.g. www.example.com)
      --user-agent string   set custom user-agent
  -v, --version             version for get-old-tweets
      --within string       find tweets nearby a certain location (e.g. 1km)
```

## Example

Get tweets by username:

```sh
get-old-tweets --from 'barackobama' -o out.csv
```

Get tweets by query search:

```sh
get-old-tweets -q 'europe refugees' -o out.csv
```

Get tweets by username and bound dates:

```sh
get-old-tweets --from 'barackobama' --since 2015-09-10 --until 2015-09-12 -o out.csv
```

Get top tweets by username:

```sh
get-old-tweets --from 'barackobama' --top -o out.csv
```

## Output CSV schema

- `id` (int)
- `username` (str)
- `created_at` (datetime)
- `full_text` (str)
- `retweet_count` (int)
- `favorite_count` (int)
- `reply_count` (int)
- `quote_count` (int)
- `latitude` (float)
- `longitude` (float)
- `lang` (str)
- `source` (str)

## Build

```sh
make build
```

## Test

```sh
make test
```
