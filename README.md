# go-graphite

![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/cybergarage/go-graphite)
[![test](https://github.com/cybergarage/go-graphite/actions/workflows/make.yml/badge.svg)](https://github.com/cybergarage/go-graphite/actions/workflows/make.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/cybergarage/go-graphite.svg)](https://pkg.go.dev/github.com/cybergarage/go-graphite)
[![Go Report Card](https://img.shields.io/badge/go%20report-A%2B-brightgreen)](https://goreportcard.com/report/github.com/cybergarage/go-graphite)
[![codecov](https://codecov.io/gh/cybergarage/go-graphite/branch/main/graph/badge.svg?token=C3Q82XPE44)](https://codecov.io/gh/cybergarage/go-graphite)

The `go-graphite` handles [the Render](https://graphite.readthedocs.io/en/latest/render_api.html) and [Carbon (feeding)](http://graphite.readthedocs.io/en/latest/feeding-carbon.html) APIs of [Graphite](https://graphiteapp.org/) so that all developers can develop [Graphite](https://graphiteapp.org/)-compatible servers easily.

![](doc/img/framework.png)

The `go-graphite` go-graphite can collect monitoring data from any products that support the [Graphite](https://graphiteapp.org/) interface ([the Carbon API](http://graphite.readthedocs.io/en/latest/feeding-carbon.html)), and allowing the use of useful visualisation tools such as [Graphana](https://grafana.com/) through [the Render API](https://graphite.readthedocs.io/en/latest/render_api.html).


## Table of Contents

- [Getting Started](doc/server_impl.md)

## References

- [The Render URL API](https://graphite.readthedocs.io/en/latest/render_api.html)
- [Feeding In Your Data](http://graphite.readthedocs.io/en/latest/feeding-carbon.html)
