# backlog-bulk-issue-registration-cli

## Description
A command line tool for bulk registering of Backlog issue.

## Installation
This package can be installed with the go get command:

```
$ go get github.com/vvatanabe/backlog-bulk-issue-registration-cli
```

Built binaries are available on Github releases:  
https://github.com/vvatanabe/backlog-bulk-issue-registration-cli/releases

## Usage

```
$ backlog-bulk-issue-registration [options]
```

## Options
```
--host value     (Required) backlog host name. Ex: xxx.backlog.jp, xxx.backlog.com [$BACKLOG_HOST]
--project value  (Required) backlog project key. [$BACKLOG_PROJECT_KEY]
--key value      (Required) backlog api key. [$BACKLOG_API_KEY]
--file value     import file path. [$BACKLOG_IMPORT_FILE]
--lang value     language setting. (ja or en) (default: "ja") [$BACKLOG_LANG]
--progress       show progress bar
--help, -h       show help
--version, -v    print the version
```

### Backlog API Key
API Key is necessary because this CLI depends on Backlog API v2.  
https://support.backlog.com/hc/en-us/articles/115015420567-API-Settings

## Friends
[backlog-bulk-issue-registration-gas](https://github.com/nulab/backlog-bulk-issue-registration-gas)

## Bugs and Feedback
For bugs, questions and discussions please use the Github Issues.

## License
[MIT License](http://www.opensource.org/licenses/mit-license.php)

## Author
[vvatanabe](https://github.com/vvatanabe)