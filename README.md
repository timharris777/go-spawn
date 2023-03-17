# go-spawn <!-- omit from toc -->

![Under Development](https://img.shields.io/badge/Under%20Development-gray?style=flat)

A cli tool written in go for project templating, scaffolding, and text-replacement
* [Install](#install)
  * [Mac](#mac)
    * [Intel](#intel)
    * [Arm](#arm)
* [Linux](#linux)
    * [x64](#x64)
    * [Arm](#arm-1)
* [Windows](#windows)
* [Use](#use)
* [Examples](#examples)
  * [A basic use case](#a-basic-use-case)
  * [Render by piping template and providing input directly](#render-by-piping-template-and-providing-input-directly)
  * [Render by piping input and providing template directly](#render-by-piping-input-and-providing-template-directly)
  * [Render a template directory](#render-a-template-directory)

## Install

### Mac

#### Intel
```sh
wget https://github.com/timharris777/go-spawn/releases/latest/download/go-spawn-darwin-amd64 -O /usr/local/bin/go-spawn && chmod +x /usr/local/bin/go-spawn
```
#### Arm
```sh
wget https://github.com/timharris777/go-spawn/releases/latest/download/go-spawn-darwin-arm64 -O /usr/local/bin/go-spawn && chmod +x /usr/local/bin/go-spawn
```

## Linux

#### x64
```sh
wget https://github.com/timharris777/go-spawn/releases/latest/download/go-spawn-linux-amd64 -O /usr/local/bin/go-spawn && chmod +x /usr/local/bin/go-spawn
```
#### Arm
```sh
wget https://github.com/timharris777/go-spawn/releases/latest/download/go-spawn-linux-arm64 -O /usr/local/bin/go-spawn && chmod +x /usr/local/bin/go-spawn
```

## Windows

Go to `https://github.com/timharris777/go-spawn/releases/latest` and download the appropriate `exe` file.

## Use

```sh
A cli tool written in go for project templating, scaffolding, and text-replacement

Usage:
  go-spawn [flags]

Flags:
  -d, --debug              folder to output rendered templates
  -h, --help               help for go-spawn
  -i, --input string       provide a yaml file that has inputs for templating
      --inputFromPipe      provide input from pipe
  -o, --output string      folder to output rendered templates
  -t, --template string    path to template file or folder. Folder requires --output option
      --templateFromPipe   provide template from pipe
```

## Examples

### A basic use case
```sh
Coming soon...
```
### Render by piping template and providing input directly
```sh
Coming soon...
```
### Render by piping input and providing template directly
```sh
Coming soon...
```
### Render a template directory
```sh
Coming soon...
```
