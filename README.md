# go-drcom

![](https://img.shields.io/github/stars/wuxiaobai24/go-drcom?style=social)
![](https://img.shields.io/github/downloads/wuxiaobai24/go-drcom/total)
![](https://img.shields.io/github/license/wuxiaobai24/go-drcom)



A simple and easy easy-to-use command line app to login drcom (in szu).

## Features

- [x] Basic Login
- [x] Loop Mode
- [x] Daemon (only support in Unix* system)
- [ ] More friendly logging
- [ ] Using configuration file
- [ ] Unit testing
- [ ] Support systemd


## Install

go-drcom is a command line app, so you can get this app in github release easily.

## Compile from source

1. Install Golang, you can check [golang doc about install](https://golang.org/doc/install)
2. Clone this repository: `git clone https://github.com/wuxiaobai24/go-drcom.git`.
3. build or install drcom: `go build` or `go install`.

## Usage

### Login in once

Threre are three way to set your username and password:

1. Input in terminal.
2. Using `--username`/`-u` and `--password`/`-p` options.
3. Set in environment variables(`GO_DRCOM_USERNAME` and `GO_DRCOM_PASSWORD`)

```bash
# 1. Input in terminal.
go-drcom 
# 2. Using `--username`/`-u` and `--password`/`-p` options.
go-drcom -u <username> -p <password>
# 3. Set in environment variables(`GO_DRCOM_USERNAME` and `GO_DRCOM_PASSWORD`)
GO_DRCOM_USERNAME=<username> GO_DRCOM_PASSWORD=<password> go-drcom 
```

### Specify Login Url

If you are in the dormitory area, may be you can set login url as http://172.30.255.2/a30.htm . You can use `--loginUrl` or `-l` to specify login url.

```bash
go-drcom -l http://172.30.255.2/a30.htm
```

### Login in loop

In loop mode, drcom will check network every <gap> seconds, if network is logout, it will login. You can use `--loop` or `--isLoop` enter loop mode, and meantime you can set gap by `--gap` or `-g` and specify testing url by `-t` and `--testingUrl`.

```bash
go-drcom --loop -g 10
```

### Daemon

In Unix* system, go-drcom can beacome a daemon by using `-d` or `--isDaemon` option. 

**Warning**: go-drcom must be in loop mode, in other words, you should use `--loop` option when you use `-d` option.

```bash
go-drcom --loop -d
```

### Help

More detail about usage, you can see the help output:

```bash
NAME:
   go-drcom - A simple and easy easy-to-use command line app to login drcom (in szu).

USAGE:
   drcom [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --gap value, -g value         Time gap for check network (default: 60)
   --isDaemon, --daemon, -d      Whether daemon (default: false)
   --isLoop, --loop              Login Loop (default: false)
   --loginUrl value, -l value    Login Url (default: "https://drcom.szu.edu.cn")
   --password value, -p value    Password for drcom [$GO_DRCOM_PASSWORD]
   --testingUrl value, -t value  Testing Url (default: "https://www.baidu.com")
   --username value, -u value    Username for drcom [$GO_DRCOM_USERNAME]
   --help, -h                    show help (default: false)
```

## LICENSE

[LPGL-2.1](https://github.com/wuxiaobai24/go-drcom/LICENSE)