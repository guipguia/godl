# Go download (godl)

Go download (godl) is a command-line interface (CLI) tool written in Go that allows you to download files from the internet with progress monitoring and support for resumable downloads.

## Features

- Download files from the internet with progress monitoring
- Resumable downloads to continue interrupted downloads
- Concurrent downloading using multiple goroutines for faster downloads

## Installation

To install godl, you need to have Go installed on your system. Then, you can use the following command to install godl cli:

```shell
go install github.com/guipguia/godl
```

## Usage

```shell
godl [options] <url> 
```

Replace <url> with the URL of the file you want to download. You can also specify the following options:

-o, --output <filename>: Specify the output filename (default: same as the filename in the URL)
-c, --concurrency <value>: Set the number of concurrent downloads (default: 4)
-r, --resume: Enable resumable downloads to continue interrupted downloads
-h, --help: Display help information
