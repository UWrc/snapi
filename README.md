# snapi

Automation for the Service Now API (i.e., UW Connect). This is a proof-of-concept tool for automating our interaction with UW Connect. The first use case was to run a command on the infrastructure then automatically capture STDOUT and STDERR in the Work Notes.

Additional details about the Service Now API can be found [here](https://uwconnect.uw.edu/kb_view.do?sysparm_article=KB0025022).

## API Keys

Create a `.snapi` file in the same location as the executable. This is a dot env file with the `SNAPI_USERNAME` and `SNAPI_PASSWORD` variables. The credentials can be found in the team shared password manager.

## Usage

Run `snapi -h` for additional information.

```bash
$ snapi -h
A command line tool to interact with the ServiceNow API.

Usage:
  snapi [flags] [command]

Flags:
  -c, --configuration-item string   Configuration item (required). (default "hyak")
  -h, --help                        help for snapi
  -k, --key string                  config file
  -r, --record string               Service Now record number (required).
  -s, --state string                The state of the record. Valid values are (o)pen or (r)esolved. (default "open")
$ 
```

## Future Stuff

* Add a `--dry-run` option to the `snapi` command to test without making changes.
* There are probably some other bugs and features to poke at.
* Address some billing automation?
