# GoBackend

[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://go.dev)
[![Github Workflow Status](https://img.shields.io/github/actions/workflow/status/windvalley/gobackend/gobackendci.yaml)](https://github.com/windvalley/gobackend/actions/workflows/gobackendci.yaml)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=windvalley_gobackend&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=windvalley_gobackend)
[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=windvalley_gobackend&metric=ncloc)](https://sonarcloud.io/summary/new_code?id=windvalley_gobackend)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=windvalley_gobackend&metric=bugs)](https://sonarcloud.io/summary/new_code?id=windvalley_gobackend)
[![License](https://img.shields.io/github/license/windvalley/gobackend)](License) <br>
![Page Views](https://views.whatilearened.today/views/github/windvalley/gobackend.svg)
[![Traffic Clones Total](https://img.shields.io/endpoint?url=https%3A%2F%2Fapi.sre.im%2Fv1%2Fgithub%2Ftraffic%2Fclones%2Ftotal%3Fgit_user%3Dwindvalley%26git_repo%3Dgobackend%26type%3Dcount%26label%3Dclones-total)](https://github.com/windvalley/traffic-clones-api)
[![Traffic Clones Uniques](https://img.shields.io/endpoint?url=https%3A%2F%2Fapi.sre.im%2Fv1%2Fgithub%2Ftraffic%2Fclones%2Ftotal%3Fgit_user%3Dwindvalley%26git_repo%3Dgobackend%26type%3Duniques%26label%3Dclones-uniques)](https://github.com/windvalley/traffic-clones-api)

Enterprise-grade Go web backend scaffolding based on
[cobra][1], [viper][2], [pflag][3], [zap][4], [gorm][5], [gin][6].

[1]: https://github.com/spf13/cobra
[2]: https://github.com/spf13/viper
[3]: https://github.com/spf13/pflag
[4]: https://github.com/uber-go/zap
[5]: https://github.com/go-gorm/gorm
[6]: https://github.com/gin-gonic/gin

## Features

- Project directories layout follows [golang-standards/project-layout](https://github.com/golang-standards/project-layout)

- The configuration items in the configuration file are unified with the command line parameters, and the command line parameters take precedence(Realized by [cobra][1], [viper][2] and [pflag][3]).

- Powerful logger built with [zap][4] and [lumberjack](https://github.com/natefinch/lumberjack), supports color, function caller, hooks, multi-outputs, rotation, etc.

- Integrated many useful middlewares, which can be flexibly configured in the configuration file.

- One server process can enable both http and https services at the same time.

- Server graceful shutdown.

- Custom your own environment variable to specify run environments(dev/test/prod).

- Supports operation logging and exposes related restful apis.

- Auto generate error code documentation file and necessary error code source files.

- Supports application process lock.

- Supports cloud native deployment.

- Use `Makefile` to manage the project efficiently and conveniently.

- Testable and maintainable codes.

## Quik Start

```sh
$ git clone --depth 1 https://github.com/windvalley/gobackend.git

$ cd gobackend
```

Edit `configs/dev.gobackend-apiserver.yaml`:

```yaml
# MySQL
mysql:
  # Default: 127.0.0.1:3306
  host: 127.0.0.1:3306
  # Default: ""
  username: "root"
  # Default: ""
  password: "123456"
  # Default: ""
  database: "gobackend"
```

Create database:

```sh
$ mysql -h 127.0.0.1 -P 3306 -uroot -p'123456' -e "create database gobackend;"
```

Run in development environment:

```sh
$ make run.dev
```

<img width="772" alt="run_dev" src="https://user-images.githubusercontent.com/6139938/144012376-df174b5e-0c5a-4318-817e-7d9b30e4f5cd.png">

## Makefile

```sh
$ make

Usage: make [TARGETS] [OPTIONS]

Targets:

   all               Make gen, lint, cover, build
   run.dev           Run in development mode.
   run.test          Run in test mode.
   build             Compile packages and dependencies to generate binary file for current platform.
   build.multiarch   Build for multiple platforms. See option PLATFORMS.
   image             Build docker images for host arch.
   push              Build docker images for host arch and push images to registry.
   lint              Check syntax and style of Go source code.
   test              Run unit test.
   cover             Run unit test and get test coverage.
   gen               Generate necessary source code files and doc files.
   clean             Remove all files that are created by building.
   help              Show this help.

Options:

   BINS        The binaries to build. Default is all commands in cmd/.
               This option is available for: make build/build.multiarch
               Example: make build BINS="apiserver otherbinary"
   IMAGES      Docker images to build. Default is all commands in cmd/.
               This option is available when using: make image/image.multiarch.
               Example: make image.multiarch IMAGES="apiserver otherbinary"
   PLATFORMS   The multiple platforms to build.
               Default is 'darwin_amd64 darwin_arm64 linux_amd64 linux_arm64 windows_amd64'.
               This option is available when using: make build.multiarch.
               Example: make build.multiarch PLATFORMS="linux_amd64"
```

## License

This project is under the MIT License.
See the [LICENSE](LICENSE) file for the full license text.
