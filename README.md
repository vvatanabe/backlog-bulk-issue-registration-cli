# backlog-bulk-issue-registration-cli [![Build Status](https://travis-ci.org/vvatanabe/backlog-bulk-issue-registration-cli.svg?branch=master)](https://travis-ci.org/vvatanabe/backlog-bulk-issue-registration-cli) [![Coverage Status](https://coveralls.io/repos/github/vvatanabe/backlog-bulk-issue-registration-cli/badge.svg?branch=master)](https://coveralls.io/github/vvatanabe/backlog-bulk-issue-registration-cli?branch=master)

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
$ bbir [options] FILE_PATH
```

If you do not specify a FILE_PATH, it will listen on standard input.

## Options
```
--host, -H value     (Required) backlog host name. Ex: xxx.backlog.jp [$BACKLOG_HOST]
--project, -P value  (Required) backlog project key.                  [$BACKLOG_PROJECT_KEY]
--key, -K value      (Required) backlog api key.                      [$BACKLOG_API_KEY]
--lang, -l value     language setting. (ja or en) (default: "en")     [$BACKLOG_LANG]
--progress, -p       show progress bar
--help, -h           show help
--version, -v        print the version
```

### Backlog API Key
API Key is necessary because this CLI depends on Backlog API v2.  
https://support.backlog.com/hc/en-us/articles/115015420567-API-Settings

## Example
From file:
```
$ bbir --host="xxx.backlog.jp" \
    --project="yourProjectKey" \
    --key="yourAPIKey" \
    ./testdata/example.csv
```

From standard input:
```
$ cat ./testdata/example.csv | bbir \
    --host="xxx.backlog.jp" \
    --project="yourProjectKey" \
    --key="yourAPIKey" \
```

## Support Fields
- Summary (Required)
- Description
- StartDate
- DueDate
- EstimatedHours
- ActualHours
- IssueTypeName (Required)
- Category
- Version
- Milestone
- Priority
- Assignee
- ParentIssue
- CustomFields

## CSV Specification
Please refer to the example below:
- [example.csv](https://github.com/vvatanabe/backlog-bulk-issue-registration-cli/blob/master/testdata/example.csv)
- [example_ja.csv](https://github.com/vvatanabe/backlog-bulk-issue-registration-cli/blob/master/testdata/example_ja.csv)

If you put an asterisk in parent issue, the above issue becomes the parent issue:
```
Header: Summary (Required), ... , ParentIssue
Line 1: Summary1          , ... ,   (This line is parent issue of line 2 and 3.)
Line 2: Summary1-1        , ... , * (This line is child issue of line 1.)
Line 3: Summary1-2        , ... , * (This line is child issue of line 1.)
```

## Friends
[backlog-bulk-issue-registration-gas](https://github.com/nulab/backlog-bulk-issue-registration-gas)

## Bugs and Feedback
For bugs, questions and discussions please use the Github Issues.

## License
[MIT License](http://www.opensource.org/licenses/mit-license.php)

## Author
[vvatanabe](https://github.com/vvatanabe)