## tfrefactor mv

Move element to a different file

### Synopsis

Move element to a different file.
Can be var, data, resource, output, local, or module.

Arguments:
  ADDRESS     The address (e.g. var.a, data.vpc.default, aws_vpc.default).
  TO_FILE     File to move to.


```
tfrefactor mv <ADDRESS> <TO_FILE> [flags]
```

### Options

```
  -c, --config string   Path of terraform to modify, defaults to current. (default "-")
  -f, --force           Skip interactive approval of update before applying
  -h, --help            help for mv
```

### SEE ALSO

* [tfrefactor](tfrefactor.md)	 - Automated refactoring for Terraform

###### Auto generated by spf13/cobra on 20-Jan-2022
