# gobackend

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

- Powerful logger built by [zap][4] and [lumberjack](https://github.com/natefinch/lumberjack), supports color, function caller, hooks, multi-outputs, rotation, etc.

- Integrated many useful middlewares, which can be flexibly configured in the configuration file.

- One server process can enable both http and https services at the same time.

- Server graceful shutdown.

- Generate error code documentation file and necessary error code source files.

- Use `Makefile` to manage the project efficiently and conveniently.

- Testable and maintainable codes.

## Quik Start

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
mysql -h 127.0.0.1 -P 3306 -uroot -p'123456' -e "create database gobackend;"
```

Project run:

```sh
./scripts/run_dev.sh
```

<img width="599" alt="rundev-shot" src="https://user-images.githubusercontent.com/6139938/141258757-b994bc59-7eee-462e-91a1-ece516035f8a.png">
