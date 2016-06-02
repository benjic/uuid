uuid
====

> You're not a beautiful and unique snowflake, but your identifiers can be.

[![Build Status](https://travis-ci.org/benjic/uuid.svg?branch=master)](https://travis-ci.org/benjic/uuid) [![codecov](https://codecov.io/gh/benjic/uuid/branch/master/graph/badge.svg)](https://codecov.io/gh/benjic/uuid) [![Go Report Card](https://goreportcard.com/badge/github.com/benjic/uuid)](https://goreportcard.com/report/github.com/benjic/uuid)

As you may of guessed this library provides [RFC4122][spec] compliant
universally unique identifiers. 

Goals
-----

- [ ] Fast
  - Benchmark should prove this library can supply a bunch of identifiers
    really, really quickly.
- [ ] Simple
  - If a consumer is fine with sane defaults consumption of library should be a
    simple factory function.
- [ ] Configurable
  - If sane defaults are not ideal, a consumer should be able to configure the
    library to suit their needs.

[spec]: https://tools.ietf.org/html/rfc4122
