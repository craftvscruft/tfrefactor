# tfrefactor
![Version](https://img.shields.io/badge/version-0.0.1-blue.svg?cacheSeconds=2592000)
[![Documentation](https://img.shields.io/badge/documentation-yes-brightgreen.svg)](https://github.com/craftvscruft/tfrefactor/blob/main/docs/cli/tfrefactor.md)
![Tests](https://github.com/craftvscruft/tfrefactor/actions/workflows/test.yml/badge.svg?branch=main)
[![License: MPL2](https://img.shields.io/github/license/raymyers/tfrefactor)](https://github.com/craftvscruft/tfrefactor/blob/main/LICENSE)

> Automated refactoring for [Terraform](https://terraform.io/).

Currently supports:

* Rename local / var / data / resource across all files in a config
* Diff preview of changes
* Adding `moved` blocks for resource renames to avoid `state mv` in Terraform 1.1

See [refactor.tf](https://refactor.tf/refactor/2021/08/26/todo.html) for more refactoring recipes.

## Install

Requires [Go 1.17](https://go.dev/doc/install)

```sh
git clone git@github.com:craftvscruft/tfrefactor.git
cd tfrefactor
make

# Ensure ~/.local/bin is in your $PATH or copy to a directory that is.
cp bin/tfrefactor ~/.local/bin
```

## Usage

Rename a var from `acct_id` to `account_id` in the current directory
```
tfrefactor rename var.acct_id var.account_id
```

Display CLI help with usage information.

```sh
tfrefactor
```

You can also run without installing
```sh
make

./bin/tfrefactor
```

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
