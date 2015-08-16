# Deis Client

`deis` is a command line utility used to interact with the [Deis](http://deis.io) open source PaaS.

Please add any [issues](https://github.com/deis/deis/issues) you find with this software to
the [Deis Project](https://github.com/deis/deis).

## Installation

Currently the only way to use the go version of the `deis` client is to build it yourself.
To build the `deis` client, you need to have [go](https://golang.org/), [Godep](https://github.com/tools/godep),
 and [make](https://www.gnu.org/software/make/) installed. Then run `make build`.

## Usage

Running `deis help` will give you a up to date list of `deis` commands.
To learn more about a command run `deis help <command>`.

## Windows Support

`deis` has experimental support for Windows. To build deis for Windows, you need to install
[go](https://golang.org/) and [Godep](https://github.com/tools/godep). Then run the `make.bat` script.

## License

Copyright 2015, Engine Yard, Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
