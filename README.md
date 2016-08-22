# Better cd
better cd

## Features
- Fast directory jumping

## Usage
All you need is `cd path` the same as `cd`.

If path exists in your current path, do the same as build in `cd`.

If path don't exist, bettercd will find a path from log and change that directory.

## Installation
#### Requirement
- peco
- go lang compilar

### how to install
Firstly, edit `Makefile` line 10 if you use bash.

`make`→`make install`

## Config
There are 4 paramters in `./bettercd/bettercd.sh`.

- BCD_DIR : bettercd directory (default:~/.bettercd)
- BCD_FILTER_TYPE : filter application (default:peco)
- BCD_LOG_DIR : Log path (default:$BCD_DIR/directory.log)
- BCD_LOG_MAX : max number of log (default:200)

## Others
#### License
MIT

#### References

- slice del func implementation : [http://jxck.hatenablog.com/entry/golang-slice-internals](http://jxck.hatenablog.com/entry/golang-slice-internals)
- line counter func implementation : [http://stackoverflow.com/questions/24562942/golang-how-do-i-determine-the-number-of-lines-in-a-file-efficiently](http://stackoverflow.com/questions/24562942/golang-how-do-i-determine-the-number-of-lines-in-a-file-efficiently)

