# tfrefactor
![Version](https://img.shields.io/badge/version-0.0.1-blue.svg?cacheSeconds=2592000)
[![Documentation](https://img.shields.io/badge/documentation-yes-brightgreen.svg)](https://refactor.tf)
[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://github.com/kefranabg/readme-md-generator/graphs/commit-activity)
![Tests](https://github.com/craftvscruft/tfrefactor/actions/workflows/test.yml/badge.svg?branch=main)
[![License: MPL2](https://img.shields.io/github/license/raymyers/tfrefactor)](https://github.com/craftvscruft/tfrefactor/blob/main/LICENSE)
[![Twitter: lambdapocalypse](https://img.shields.io/twitter/follow/lambdapocalypse.svg?style=social)](https://twitter.com/lambdapocalypse)

> Automated refactoring for [Terraform](https://terraform.io/).

### 🏠 [Homepage](https://github.com/craftvscruft/tfrefactor)

## Install

Requires [Go 1.17](https://go.dev/doc/install)

Substitute ~/.local/bin for a prefered dir the is in your $PATH.

```sh
git clone git@github.com:craftvscruft/tfrefactor.git
cd tfrefactor
make

cp bin/tfrefactor ~/.local/bin
```

## Usage

```sh
tfrefactor
```

Or without installing / during development
```sh
make
./bin/tfrefactor
```

The CLI help will explain usage.

## Run tests

```sh
make test
```

## Author

👤 **Ray Myers**

* YouTube: [Craft vs Cruft](https://www.youtube.com/channel/UC4nEbAo5xFsOZDk2v0RIGHA)
* Twitter: [@lambdapocalypse](https://twitter.com/lambdapocalypse)
* GitHub: [@raymyers](https://github.com/raymyers)
* LinkedIn: [@cadrlife](https://linkedin.com/in/cadrlife)

## 🤝 Contributing

Contributions, issues and feature requests are welcome!

Feel free to check [issues page](https://github.com/craftvscruft/tfrefactor/issues). You can also take a look at the [contributing guide](https://github.com/craftvscruft/tfrefactor/blob/main/CONTRIBUTING.md).

## Show your support

Give a ⭐️ if this project helped you!

[![support us](https://img.shields.io/badge/become-a%20patreon%20us-orange.svg?cacheSeconds=2592000)](https://www.patreon.com/craftvscruft)

## Acknowledgements

Built on [hclwrite](https://github.com/hashicorp/hcl), a component of HashiCorp Configuration Language (HCL).

Inspiration and test helper code from [hcledit](https://github.com/minamijoyo/hcledit) by Masayuki Morita.

## 📝 License

Copyright © 2022 [Ray Myers](https://github.com/raymyers).

This project is [MPL2](https://github.com/craftvscruft/tfrefactor/blob/main/LICENSE) licensed.
