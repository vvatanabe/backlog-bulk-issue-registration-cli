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
$ bbir [options]
```

## Options
```
--host, -H value     (Required) backlog host name. Ex: xxx.backlog.jp [$BACKLOG_HOST]
--project, -p value  (Required) backlog project key.                  [$BACKLOG_PROJECT_KEY]
--key, -k value      (Required) backlog api key.                      [$BACKLOG_API_KEY]
--file, -f value     import file path.                                [$BACKLOG_IMPORT_FILE]
--lang, -l value     language setting. (ja or en) (default: "ja")     [$BACKLOG_LANG]
--progress, -P       show progress bar
--help, -h           show help
--version, -v        print the version
```

### Backlog API Key
API Key is necessary because this CLI depends on Backlog API v2.  
https://support.backlog.com/hc/en-us/articles/115015420567-API-Settings

## Example
From file:
```
$ bbir --file="./testdata/example.csv" \
    --host="xxx.backlog.jp" \
    --project="yourProjectKey" \
    --key="yourAPIKey" \
```

From standard input:
```
$ cat ./testdata/example.csv | bbir \
    --host="xxx.backlog.jp" \
    --project="yourProjectKey" \
    --key="yourAPIKey" \
```

## Friends
[backlog-bulk-issue-registration-gas](https://github.com/nulab/backlog-bulk-issue-registration-gas)

## Bugs and Feedback
For bugs, questions and discussions please use the Github Issues.

## License
[MIT License](http://www.opensource.org/licenses/mit-license.php)

## Author
[vvatanabe](https://github.com/vvatanabe)