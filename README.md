slackbotwit
===========

A slackbot hook with WIT.ai for NLP. I use this package "github.com/abourget/slick"
to take advantage of pluggable architecture.

Build Status: [![Build Status](https://travis-ci.org/lenfree/slackbotwit.svg?branch=master)](https://travis-ci.org/lenfree/slackbotwit)

Status: Work In Progress

[Binary Releases](https://github.com/lenfree/slackbotwit/releases)

### Getting started:
--------------------

### Install required packages:
```
$ make install
```

### Run test:
```
$ make script
```

### Start server with binary:
```
1. Download binary release
2. cp env.example .env
3. ./bbot
```

### Start server without binary:
```
1. cp env.example .env
2. make run
```

## Build binary for Pi
```
$ make buildpi VERSION=0.0.25
```

## Sync to Pi
```
$ make shippi VERSION=0.0.25
```

### Contributing:
```
1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
```
