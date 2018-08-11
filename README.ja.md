# backlog-bulk-issue-registration-cli [![Build Status](https://travis-ci.org/vvatanabe/backlog-bulk-issue-registration-cli.svg?branch=master)](https://travis-ci.org/vvatanabe/backlog-bulk-issue-registration-cli) [![Coverage Status](https://coveralls.io/repos/github/vvatanabe/backlog-bulk-issue-registration-cli/badge.svg?branch=master)](https://coveralls.io/github/vvatanabe/backlog-bulk-issue-registration-cli?branch=master)

## Description
これは Backlog の課題を一括登録するためのコマンドラインツールです。

## Installation
This package can be installed with the go get command:

```
$ go get github.com/vvatanabe/backlog-bulk-issue-registration-cli
```

Built binaries are available on Github releases:  
https://github.com/vvatanabe/backlog-bulk-issue-registration-cli/releases

## 使い方

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
--check, -c          check mode (validation only)
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
- 件名 (必須)
- 詳細
- 開始日
- 期限日
- 予定時間
- 実績時間
- 種別名 (必須)
- カテゴリ名
- バージョン名
- マイルストーン名
- 優先度名
- 担当者名
- 親課題
- カスタム属性

## 入力項目について (CSV)
- テンプレートを元に、CSVファイルの内容を登録したい情報に書き換えてください。
- １行目のヘッダーは、登録処理に必要な情報を含んでいるため、削除したり内容を変更しないでください。
- 「件名」「種別」は登録するための必須項目です。

### 親課題
- 親課題が既に存在する場合は課題キーを指定してください。
- CSVファイル内の課題を指定する場合は課題キーがまだ存在しないので、`*` を入力することで直近の親課題を指定していない課題を親課題とします。

### カスタム属性
- プロジェクトにカスタム属性が登録されている場合、ヘッダーにカスタム属性の名前を指定るすることで登録できるようになります。
- 複数選択可能なカスタム属性は追加できません。
- カスタム属性はプレミアムプラン以上からご利用できます。


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