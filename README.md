# Go MYSQL Dump
Create MYSQL dumps on Symfony projects in Golang, by parsing parameters.yml file to get database connection parameters.

### Usage Example
```console
$ go build export.go
$ ./export /path/to/parameters.yml /path/to/dumps [optional_path_to_sock]
```
