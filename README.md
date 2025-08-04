<p align="center">
  <img src="img/logo.png" alt="snapi" style="width:300px;"/>
  <br />
  <a href="https://github.com/UWrc/snapi/actions/workflows/go.yml"><img src="https://github.com/UWrc/snapi/actions/workflows/go.yml/badge.svg?branch=main"></a>
</p>

# snapi

Automation for the Service Now API (i.e., UW Connect). This is a proof-of-concept tool for automating our interaction with UW Connect. The first use case was to run a command on the infrastructure then automatically capture STDOUT and STDERR in the Work Notes.

Additional details about the Service Now API can be found [here](https://uwconnect.uw.edu/kb_view.do?sysparm_article=KB0025022) for changes (e.g., REQ, CHG, INC) and [here](https://uw.service-now.com/now/nav/ui/classic/params/target/kb_view.do%3Fsysparm_article%3DKB0032495) for billing.

## API Keys

Create a `.snapi` file in the same location as the executable. This is a dot env file with the `SNAPI_USERNAME` and `SNAPI_PASSWORD` variables. The credentials can be found in the team shared password manager.

## Usage

Run `snapi -h` for additional information. `snapi` is available on `klone-head01`.

```bash
$ snapi -h
snapi v0.0.3: A command line tool to interact with the ServiceNow API.

Usage:
  snapi [flags] [command]

Flags:
  -a, --assigned-to string          A single netID or email address for the primary contact working on the record.
  -c, --configuration-item string   Configuration item (required). (default "hyak")
  -h, --help                        help for snapi
  -k, --key string                  API key file.
  -n, --note-list string            A comma-separated list of email addresses to add to the work note watch list for this record. This is for all internal (i.e., non-customer) facing communications and notes.
  -r, --record string               Service Now record number (required). Only REQs, CHGs, and INCs supported.
  -s, --state string                The state of the record. Valid values are (o)pen or (r)esolved. (default "open")
  -w, --watch-list string           A comma-separated list of email addresses to add to the watch list for this record. This is for all customer facing communications.
$ 
```

## Future Stuff

* Add a `--dry-run` option to the `snapi` command to test without making changes.
* There are probably some other bugs and features to poke at.
* Address some billing automation?
