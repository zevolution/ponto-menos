<span id="top"></span>
<p align="center">
<img src="https://github.com/spring-projects/spring-boot/assets/36534847/b6a65f6b-39db-4b9a-b509-0553d736811e" alt="PontoMenos-Logo">

## Table of contents
- [About the project](#about-the-project)
- [Description](#description)
- [Built with](#built-with)
- [Installation](#installation)
- [Requirements to run](#requirements-to-run)
- [Modules](#modules)
- [License](#license)

## About the project
This project is a part of my private project portfolio originally implemented in Java. The primary motivation for its migration is to enhance my proficiency with the Go programming language.
<p align="right">(<a href="#top">back to top</a>)</p>

## Description
Are you tired of forgetting to clock in or simply are very lazy like me and find it cumbersome to navigate to the clock-in page every time? Stop worrying! The Ponto-Menos app offers a seamless solution for punctuality on [Ponto Mais](https://pontomais.com.br), the brazilian platform responsible for managing employee punch clocks, now under the [VR](https://www.vr.com.br) umbrella. Simplify your workday routine with this user-friendly tool and never miss the "time" again!
<p align="right">(<a href="#top">back to top</a>)</p>

## Built with
* [Golang](https://go.dev)
* [Serverless Framework](https://www.serverless.com)
* [Cobra](https://github.com/spf13/cobra)
* [AWS Lambda Go](https://github.com/aws/aws-lambda-go)
* [Zap](https://github.com/uber-go/zap)
* [Viper](https://github.com/spf13/viper)
* [LumberJack](https://github.com/natefinch/lumberjack)
<p align="right">(<a href="#top">back to top</a>)</p>

## Installation

To clone and run this application, you'll need Git installed on your computer(or no, if you want to download **.zip**). From your command line:
```bash
# Git CLI
git clone https://github.com/zevolution/ponto-menos.git

# Github CLI
gh repo clone zevolution/ponto-menos
```
<p align="right">(<a href="#top">back to top</a>)</p>

## Requirements to run
* If you use Windows OS, is strongly recommended that you use WLS2 to perform all operations.
* [Golang 1.19](https://tip.golang.org/doc/go1.19)
<p align="right">(<a href="#top">back to top</a>)</p>

## Build
This project uses Makefile, to view build options using `make help`:
```bash
build                          Build all binaries
build-cli                      Build CLI binary (e.g. 'make os=darwin arch=arm64 build-cli' or just 'make build-cli' for linux)
build-lambda                   Build lambda binaries using param 'name' (e.g. make name=lambda-name build-lambda)
zip-lambda                     Zip lambda binaries
clean                          Remove previous build
help                           Display available commands
```
<p align="right">(<a href="#top">back to top</a>)</p>

## Modules
This project features two modules: [Punch-Clock Schedule](cmd/punchclockschedule/README.md) and [CLI](cmd/cli/README.md). For in-depth details about each module, please explore the following links.

<details>
    <summary>
        <strong>
            <span><a href="https://github.com/zevolution/ponto-menos/tree/main/cmd/punchclockschedule/README.md">⏱️ Punch-Clock Schedule</a></span>
        </strong>
    </summary>
    <img src="https://github.com/spring-projects/spring-boot/assets/36534847/a59e1e3e-7eae-4e1d-9aec-53e128b10992"/>
</details>

<details>
    <summary>
        <strong>
            <span><a href="https://github.com/zevolution/ponto-menos/tree/main/cmd/cli/README.md">🖥️ CLI</a></span>
        </strong>
    </summary>
    <img src="https://github.com/spring-projects/spring-boot/assets/36534847/3d55b7a2-bbdc-436f-8f78-88b9197669a9"/>
</details>
<p align="right">(<a href="#top">back to top</a>)</p>

## License
[MIT](https://choosealicense.com/licenses/mit/)
