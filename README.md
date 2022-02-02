# loggy
---

[![Go Report Card](https://goreportcard.com/badge/github.com/auhau/loggy)](https://goreportcard.com/report/github.com/auhau/loggy)
[![Build Status](https://img.shields.io/endpoint.svg?url=https%3A%2F%2Factions-badge.atrox.dev%2Fauhau%2Floggy%2Fbadge&style=flat&label=build)](https://actions-badge.atrox.dev/auhau/loggy/goto)

> Swiss knife for logs

## Installation

### Automatic installation script

```shell
$ curl -sSL https://raw.githubusercontent.com/auhau/loggy/master/install.sh | bash
```

### Ubuntu / Debian / Raspbian / CentOS

`deb` and `rpm` packages are build. Get the one for your OS
from  [latest release](https://github.com/auhau/loggy/releases/latest) page and run `sudo dpkg -i <package>.deb` or
`sudo rpm -i <package>.rpm`.

### Gophers

If you've already got a Go development environment set up, you can grab it like this:

```shell
$ go get github.com/auhau/loggy
```

### Homebrew

```shell
$ brew tap auhau/loggy
$ brew install loggy
```

### Scoop

```
scoop bucket add org https://github.com/auhau/scoop.git
scoop install org/loggy
```

### Manual

Download on to your `$PATH` loggy binary for your operating system and architecture
from [latest release](https://github.com/auhau/loggy/releases/latest) page.

## Usage

See `loggy --help` for this usage:

```
By default loggy reads from STDIN or you can specify file path to read the logs from specific file.

Configuration
-------------
All options that can be passed using CLI flags can be configured using config file or environment variables.
loggy looks for a config file ".loggy.toml" in current working directory and in $HOME folder. 
Moreover you can use environment variables with "LOGGY_" prefix, where for example flag "--pattern" would be "LOGGY_PATTERN" env. variable.
The order of precedence is: $HOME config > CWD config > --config config > Env. variables > CLI flags.

Parsing pattern
---------------
The logs are parsed using parsing pattern that you have to configure in order to use filters. The lines are tokenized using space character so you have to define section of the line. Internally regex is used for parsing, but the input pattern is escaped by default for special characters so you don't have to worry about that. You define parameters using syntax "<name:type>", where name is the name of parameter that you can refer to in filters and type is predefined type used to correctly find and parse the parameter.

Supported types:
 - "string" defines string containing non-whitespace characters: [^\s]+
 - "integer" defines a integer: [0-9]+
 - "rest" collects the rest of the line: .*

Example log and bellow its parsing pattern:
[2022-09-11T15:04:22](authorization) DEBUG 200 We have received login information
[<timestamp:string>](<component:string>) <level:string> <code:integer> <message:rest>

Pattern names
-------------
In your config file you can create a [patterns] section where you can predefine your patterns using <name>="<pattern>" syntax and then use --pattern-name/-n flag to use it.

Filter
------
In order to use filter for the logs you have to define parsing pattern in which you define parameters that are extracted from the log lines. Then you can write filter expressions that will be applied on the logs. Filter has to return bool otherwise error will be shown.

loggy uses internally "govaluate" which has very rich set of C-like arithmetic/string expressions that you can use for your filters. Brief overview:
 - modifiers: + - / * & | ^ ** %% >> <<
 - comparators: > >= < <= == != =~ !~
 - logical ops: || &&
 - numeric constants, as 64-bit floating point (12345.678)
 - string constants (single quotes: 'foobar')
 - date constants (single quotes, using any permutation of RFC3339, ISO8601, ruby date, or unix date; date parsing is automatically tried with any string constant)
 - boolean constants: true false
 - parenthesis to control order of evaluation ( )
 - arrays (anything separated by , within parenthesis: (1, 2, 'foo'))
 - prefixes: ! - ~
 - ternary conditional: ? :
 - null coalescence: ??

For more details see: https://github.com/Knetic/govaluate/blob/master/MANUAL.md

Example of filter for the parsing pattern log above:
level == "DEBUG" - display only debug messages
code > 400 - display logs with code higher then 400

Keyboard shortcuts
------------------
General keys:
 - "/" for setting filter
 - "f" for toggling filter
 - "p" for setting parsing pattern input
 - "b" for scroll to bottom and follow new data
 - "h" for displaying help

Logs navigation:
 - "j", "k" or arrow keys for scrolling by one line 
 - "g", "G" to move to top / bottom
 - "Ctrl-F", "page down" to move down by one page
 - "Ctrl-B", "page up" to move up by one page

Input fields:
 - Left arrow: Move left by one character.
 - Right arrow: Move right by one character.
 - Home, Ctrl-A, Alt-a: Move to the beginning of the line.
 - End, Ctrl-E, Alt-e: Move to the end of the line.
 - Alt-left, Alt-b: Move left by one word.
 - Alt-right, Alt-f: Move right by one word.
 - Backspace: Delete the character before the cursor.
 - Delete: Delete the character after the cursor.
 - Ctrl-K: Delete from the cursor to the end of the line.
 - Ctrl-W: Delete the last word before the cursor.
 - Ctrl-U: Delete the entire line.

Usage:
  loggy [path to log file] [flags]

Flags:
  -b, --buffer-size int       number of lines that will be buffered (default 10000)
      --config string         config file (default is $HOME/.loggy.yaml)
  -h, --help                  help for loggy
  -p, --pattern string        parsing pattern see above for details
  -n, --pattern-name string   use predefined pattern in config
```

## Contribute

PRs are welcome!
If you plan to make some bigger change then please discuss it previously in Issue so there would not be any
misunderstanding and your work would get merged!

## License

[MIT](./LICENSE)