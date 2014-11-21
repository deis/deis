# Targz

Library for packaging/extracting folders in tar.gz archives.

[Documentation on godoc.org](http://godoc.org/github.com/walle/targz)

## Installation

Installing using go get is the easiest.

    go get github.com/walle/targz

## Usage

The API is really simple, there are only two methods.

* Compress
* Extract

### Create an archive containing a folder

    import "github.com/walle/targz"
    ...
    err := targz.Compress("my_folder", "my_file.tar.gz")

### Extract an archive

    import "github.com/walle/targz"
    ...
    err := targz.Extract("my_file.tar.gz", "path/to/extract/to")

## Contributing

All contributions are welcome! See [CONTRIBUTING](CONTRIBUTING.md) for more info.

## License

Licensed under MIT license. See [LICENSE](LICENSE) for more information.