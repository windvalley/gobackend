# gobackend

Enterprise-grade Go web backend scaffolding based on
[cobra][1], [viper][2], [pflag][3], [zap][4], [gorm][5], [gin][6].

[1]: https://github.com/spf13/cobra
[2]: https://github.com/spf13/viper
[3]: https://github.com/spf13/pflag
[4]: https://github.com/uber-go/zap
[5]: https://github.com/go-gorm/gorm
[6]: https://github.com/gin-gonic/gin

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
