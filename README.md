<p align="center"><a href="https://xuri.me/aurora" target="_blank" alt="Aurora Beanstalkd Console"><img width="440" src="./aurora.png" alt="aurora"></a></p>

<p align="center">
    <a href="https://travis-ci.com/xuri/aurora"><img src="https://travis-ci.com/xuri/aurora.svg?branch=master" alt="Build Status"></a>
    <a href="https://bestpractices.coreinfrastructure.org/projects/2366"><img src="https://bestpractices.coreinfrastructure.org/projects/2366/badge" alt="CII Best Practices"></a>
    <a href="https://goreportcard.com/report/github.com/xuri/aurora"><img src="https://goreportcard.com/badge/github.com/xuri/aurora" alt="Go Report Card"></a>
    <a href="https://github.com/xuri/aurora/releases"><img src="https://img.shields.io/github/downloads/xuri/aurora/total.svg" alt="Downloads"></a>
    <a href="https://github.com/xuri/aurora/blob/master/LICENSE"><img src="https://img.shields.io/github/license/mashape/apistatus.svg" alt="Licenses"></a>
    <a href="https://github.com/xuri/aurora/releases"><img src="https://img.shields.io/github/release/xuri/aurora.svg?label=Release" alt="Release"></a>
</p>

## Overview

aurora is a web-based Beanstalkd queue server console written in Go and works on macOS, Linux, and Windows machines. The main idea behind using Go for backend development is to utilize the ability of the compiler to produce zero-dependency binaries for multiple platforms. aurora was created as an attempt to build a very simple and portable application to work with a local or remote Beanstalkd server.

[See application screenshots](https://github.com/xuri/aurora/wiki)

## Features

- Cross-platform support macOS/Linux/Windows 32/64-bit
- Simple installation (distributed as a single binary)
- Zero dependencies
- Common list of servers in the config for all users + optional Basic Auth
- The full list of available tubes
- Complete statistics about jobs in tubes
- Real-time auto-update with highlighting of changed values
- You can view jobs in ready/delayed/buried states in every tube
- You can add/kick/delete jobs in every tube
- You can select multiple tubes by regExp and clear them
- You can set the statistics overview graph for every tube
- You can move jobs between tubes
- Ability to Pause tubes
- Search jobs data field
- Customizable UI (code highlighter, choose columns, edit auto refresh seconds, pause tube seconds)

## Installation

Installing aurora using [Homebrew](https://brew.sh) on macOS:

```bash
brew install aurora
```

Building aurora using Docker:

```bash
docker build -t aurora:latest .
docker run --rm --detach -p 3000:3000 aurora:latest
```

[Precompiled binaries](https://github.com/xuri/aurora/releases) for supported operating systems are available.

## Contributing

Contributions are welcome! Open a pull request to fix a bug, or open an issue to discuss a new feature or change.

## Todo

- Handle 404 error page
- Filter the tubes by name in the overview
- Logout support when Basic Auth has been enabled
- Custom job content hightlight display theme support
- Cookies control, each user can add its own personal Beanstalkd server

## Credits

- Client for beanstalkd use [beanstalkd/go-beanstalk](https://github.com/beanstalkd/go-beanstalk)
- TOML parser use [BurntSushi/toml](https://github.com/BurntSushi/toml)
- Web UI originally by [ptrofimov/beanstalk_console](https://github.com/ptrofimov/beanstalk_console)
- The logo is originally by [Ali Irawan](http://www.solusiteknologi.co.id/using-supervisord-beanstalkd-laravel/). This artwork is an adaptation

## Contributors

This project exists thanks to all the people who contribute.

[![Contributors](https://opencollective.com/aurora/contributors.svg?width=890&button=false)](https://github.com/xuri/aurora/graphs/contributors)

## Backers

Thank you to all our backers! 🙏 [Become a backer](https://opencollective.com/aurora#backer)

## Sponsors

Support this project by [becoming a sponsor](https://opencollective.com/aurora#sponsor). Your logo will show up here with a link to your website.

## Licenses

This program is under the terms of the MIT License. See [LICENSE](https://github.com/xuri/aurora/blob/master/LICENSE) for the full license text.
