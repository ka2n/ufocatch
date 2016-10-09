# ufocatch

ufocatch is a command line interface for [有報キャッチャー](ufocatch.com).

## Usage

List documents

```bash
# from edinet(with XBRL)
$ ufocatch list "4751 四半期報告書" --source=edinetx

# from tdnet(with XBRL)
$ ufocatch list "4751 決算短信" --source=tdnetx

# Show Usage
$ ufocatch --help list
```

Get resource

```bash
# XBRL as ZIP archive(default)
$ ufocatch get 'ED2014121600183' --format=xbrl

# PDF
$ ufocatch get 'ED2014121600183' --format=pdf

# Show Usage
$ ufocatch --help get

# Combination with list command
$ ufocatch list 'query' | head -n1 | ufocatch get
```

## Install

To install, use `go get`:

```bash
$ go get -d github.com/ka2n/ufocatch
```

## Contribution

1. Fork ([https://github.com/ka2n/ufocatch/fork](https://github.com/ka2n/ufocatch/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[ka2n](https://github.com/ka2n)
