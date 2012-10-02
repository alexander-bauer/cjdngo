# cjdngo

`cjdngo` is a [Go](http://golang.org) package to wrap [cjdns](https://github.com/cjdelisle/cjdns) as minimally, but usefully, as possible. For the importance of cjdns, please see [Project Meshnet](https://projectmeshnet.org).

The most fully-featured part of this package is that which wraps the JSON configuration files produced by cjdns. It provides an object, `Conf`, to allow for editing of the config file by other Go applications.

## Warning

#### Nasal Demons

`cjdngo` is still under development, and may, though some misfortune, destroy any config file that you run it on. *Please* back up your config files, and keep them safe.

#### Comment Removal

JSON does not support comments, even though `cjdns` is able to parse it properly with them present. The JSON library available in Go is not able to, however, and therefore `cjdngo` must remove all explanatory comments in your config file before acting on it. **Be warned**: Any config file that is passed through `cjdngo` **will** have all of its comments permanently removed.

## Installation

`cjdngo` is a Go package, also describable as a "cjdns library." It provides no functionality of its own. Other Go programs, however, are able to import `cjdngo` and use it to interact with cjdns. If you are developing such an application, and want to be able to import `cjdngo` just like other packages, you need to make sure that your environment is set up to be able to do so. See [here](http://golang.org/doc/code.html#tmp_79) for a much more detailed explanation.

The environment variable `$GOPATH` must point to a directory (such as `$HOME/development/go`, without the trailing `/`) containing three subdirectories: `src/`, `pkg/`, and `bin/`. Once it is doing so, you can import `cjdngo` by doing one of the following.

You can get the source locally by cloning this repository into a subdirectory of `src/`, and installing afterward, like so.

```sh
cd $GOPATH/src/
git clone git@github.com:SashaCrofter/cjdngo.git
cd cjdngo/
go install
```

Alternatively, you can use the `go` tool to fetch the package directly from GitHub.
```sh
go get github.com/SashaCrofter/cjdngo
```

If there is no output from the `go` command, then that means the installation was successful. After this point, you can use `cjdngo` by doing `import "cjdngo"` in your code.
