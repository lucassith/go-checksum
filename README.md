# go-checksum
Drag&amp;drop SHA-384 generator

## Requirements

Golang >= 1.8

## How to use

First, you need to build the binary:

```
go build main.go
```

### GUI

#### Windows & Linux only (XFCE tested)

Drag and drop files you wish to generate checksums for onto the binary/executable. Then the application will create checksums.txt file for you.

### CLI

Available commands:

`./main -h`

#### To output to the checksums.txt

```
./main {file1} {file2} {file3} ...
```

### To output to console

```
./main -c {file1} {file2} {file3} ...
```
