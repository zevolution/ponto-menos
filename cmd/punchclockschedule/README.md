<span id="top"></span>

<div align="center">

*[Back to the main README](../../README.md#modules)*

</div>

## Table of contents
- [Description](#description)
- [Build](#build)
- [Usage information](#usage-information)
- [Deploy](#deploy)
- [Disclaimer](#disclaimer)

## Description
This module was designed to offer an interactive and efficient solution for punch the clock using a cloud-native solution.
<p align="right">(<a href="#top">back to top</a>)</p>

## Build
As mentioned in [main README](../../README.md) "this project uses Makefile", therefore, we can build the application using `make clean build` to build all binaries or you can use the `make name=punchclockschedule clean build-lambda zip-lambda` to build just this app.
<p align="right">(<a href="#top">back to top</a>)</p>

## Usage information
To use this application as a cloud native solution in AWS you need make deploy using Serverless Framework, for this, you need firstly see the [disclaimer section](#disclaimer) and after this, follow the [deploy section](#deploy)
<p align="right">(<a href="#top">back to top</a>)</p>

## Deploy
To deploy this app, you'll need:
1. Install [Serverless Framework](https://www.serverless.com/framework/docs/getting-started) globally using:
```bash
npm install -g serverless
```
2. Configure your [AWS Credentials](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html). It's recommended that you use some credentials that have Administrator level privilege access, if you don't have something like this actually, you can see the following [link](https://www.serverless.com/framework/docs/providers/aws/guide/credentials/) about how to do this
3. Build application using the previous topic ["Build"](#build)
4. Make sure you are in the [punchclockschedule application directory](../punchclockschedule/)
5. Create a `.env` file following the [.env.sample](.env.sample) and configure with your data
6. Use the command `sls deploy --verbose`
7. Wait until all resources are created
8. Enjoy!
<p align="right">(<a href="#top">back to top</a>)</p>

## Disclaimer
The material embodied in this software is provided to you "as-is", we are exempt from any responsibilities related to AWS (Amazon Web Services) hereinafter determined as a cloud-provider, such as payments, fees and charges. It is the entire responsibility of the application user to deal with which events related to the cloud-provider.


<p align="right">(<a href="#top">back to top</a>)</p>