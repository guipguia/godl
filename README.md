# Go download (godl)

Go download (godl) is a command-line interface (CLI) tool written in Go that allows you to download files from the internet with progress monitoring and support for resumable downloads.

## Features

- Download files from the internet with progress monitoring
- Resumable downloads to continue interrupted downloads
- Concurrent downloading using multiple goroutines for faster downloads

## Installation

To install godl, you need to have Go installed on your system. Then, you can use the following command to install godl cli:


* Download binary from [RELEASE](https://github.com/guipguia/godl/releases)
* Add it to your PATH

## Usage

```shell
godl [options] <url> 
```

Replace <url> with the URL of the file you want to download. You can also specify the following options:

* -d --directory <path>: Specify where you want to save the file.
* -o, --output <filename>: Specify the output filename (default: same as the filename in the URL)
* -c, --concurrency <value>: Set the number of concurrent downloads (default: 4)
* -h, --help: Display help information

Also all the flags can be set as environment variables to be always a default vaule.

| Description      | Flag               | Env Var        |
|------------------|--------------------|----------------|
| Change directory | -d / --directory   | GL_DIRECTORY   |
| Change filename  | -f / --filename    | GL_FILENAME    |
| Concurrency      | -c / --concurrency | GL_CONCURRENCY |

> **WARNING**
>
> Be careful with concurrency it can actually slow down the download speed and can overwhelm your computer.


## Example

![Example of output](https://github.com/guipguia/godl/blob/main/assets/example.gif)

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvement, please open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](https://github.com/guipguia/godl/blob/main/LICENSE).