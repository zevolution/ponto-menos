<span id="top"></span>

<div align="center">

*[Back to the main README](../../README.md#modules)*

</div>

## Table of contents
- [Description](#description)
- [Build](#build)
    - [Supported Platforms](#supported-platforms)
- [Usage information](#usage-information)
    - [Logs](#logs)
    - [Env](#env)
- [Bonus](#bonus)
    - [By-step](#by-step)

## Description
This module was designed to offer an interactive and efficient solution for punch the clock on Unix-based systems, eliminating any unnecessary complexity or challenges.
<p align="right">(<a href="#top">back to top</a>)</p>

## Build
To build the CLI, use the provided Makefile with the following command: `make clean build-cli`. This command accepts two parameters called `os` and `arch`, which must be used in build for target OS correctly e.g:
```bash
make os=darwin arch=arm64 clean build-cli
```
<p align="right">(<a href="#top">back to top</a>)</p>

### Supported Platforms
If you don't know which platform is supported for build, you can use the following command:
```bash
go tool dist
```

For a better visualization of this data you can use the following command to see in columns:
```bash
go tool dist list | column -c 80 | column -t
```

Or even, if you have [jq](https://jqlang.github.io/jq/) installed on your OS, you can use the following command to get a more strucuted view:
```bash
go tool dist list -json | jq 'group_by(.GOOS) | map({ key: (.[0].GOOS), value: [.[] | .] }) | from_entries'
```
<p align="right">(<a href="#top">back to top</a>)</p>

## Usage information
To see all commands provided by Ponto-Menos CLI, you can use `./ponto-menos --help`, but in this case, currently the main used command is `time-register` responsible to punch the clock. To see the necessary params you can you `./ponto-menos time-register --help` e.g:
```bash
This command will clock in for you on the PontoMais platform using the current time

Usage:
  ponto-menos time-register [flags]

Flags:
  -h, --help              help for time-register
      --password string   Set password used in PontoMais
      --user string       Set user e-mail used in PontoMais
```

So, in this case to punch the clock use:
```bash
./ponto-menos time-register --user "your-email@domain.com" --password "yourpassword"
```

### Logs
By default log level starts in debug, the log `ponto-menos.log` is generated in the same folder as the binary is located in

### Env
To change the default behavior of the system you can use the [TOML](https://toml.io/en/) file, if you don't know this structure, think in an improved `.ini` file. To work correctly, you need to add the file `ponto-menos.toml` in the same directory where the binary is located.

You can find the sample of this file with all changeable configurations [here](ponto-menos.toml.sample)
<p align="right">(<a href="#top">back to top</a>)</p>

## Bonus
If you don't want follow the cloud solution provided [here](../punchclockschedule/README.md), you can create a schedule in your OS and call the binary to punch the clock!

Follow the diagram below:
![Cron](https://github.com/spring-projects/spring-boot/assets/36534847/302e65ea-05c7-4d67-a97c-f5def37f06b7)

### By-step
1. First of all, you need compile the binary using the command `make clean build-cli`, don't forget to enter your `os` and `arch` correctly
2. If you use:
    1. Unix-based system: the following [link]() will give you an initial idea of how use `crontab` to do this
    2. Windows: I strongly recommended that you use WLS2 to perform all operations, otherwise, you can take a look at this related content about [how to use windows task scheduler](https://learn.microsoft.com/en-us/windows/win32/taskschd/using-the-task-scheduler)
3. Enjoy!
<p align="right">(<a href="#top">back to top</a>)</p>