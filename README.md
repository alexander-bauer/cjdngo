# cjdngo

`cjdngo` is a [Go](http://golang.org) package to wrap [cjdns](https://github.com/cjdelisle/cjdns) as minimally, but usefully, as possible. For the importance of cjdns, please see [Project Meshnet](https://projectmeshnet.org).

The most fully-featured part of this package is that which wraps the JSON configuration files produced by cjdns. It provides an object, `Conf`, to allow for editing of the config file by other Go applications.
