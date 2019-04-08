# QiniuDrive

A cross-platform client for Qiniu buliding on [andlabs/ui](https://github.com/andlabs/ui) and [qiniu/api.v7](https://github.com/qiniu/api.v7). Support simple file operations `upload`, `download`, `delete`, `fetch` and extra `encrypt`, `decrypt`.

## Get started

Download the [new release](https://github.com/p1gd0g/QiniuDrive/releases).

|Platform |Login  |File |
|---|:---:|:---:|
|Linux    |![loginWindow](/images/loginWindow_linux.png)|![fileWindow](/images/fileWindow_linux.png)|
|Windows  |||
|Mac OS X |I guess|it looks similar to Linux|

To login more convieniently, excute the binary with flags in terminal, run `qiniudrive -h` to show more info:

```
Usage of qiniudrive:
  -ak string
        Your access key.
  -bk string
        Your bucket name.
  -dm string
        Your domain name. Optional for Download.
  -sk string
        Your secret key.
  -zn int
        Your bucket zone: 1-Huadong, 2-Huabei, 3-Huanan, 4-Beimei. Optional for Upload.
```

Or hardcode your user info directly.

## Build

It requires Go 1.9 or newer.

First `go get` my [fork](https://github.com/p1gd0g/ui) of [andlabs/ui](https://github.com/andlabs/ui):

```
go get -u github.com/p1gd0g/ui/...
```

Then `go get` [qiniu/api.v7](https://github.com/qiniu/api.v7):

```
go get -u github.com/qiniu/api.v7
```

Note that [andlabs/ui](https://github.com/andlabs/ui) requires:

```
- Windows: cgo, Windows Vista SP2 with Platform Update and newer
- Mac OS X: cgo, Mac OS X 10.8 and newer
- other Unixes: cgo, GTK+ 3.10 and newer
	- Debian, Ubuntu, etc.: `sudo apt-get install libgtk-3-dev`
	- Red Hat/Fedora, etc.: `sudo dnf install gtk3-devel`
```


### License

This project is licensed under the terms of the MIT license.