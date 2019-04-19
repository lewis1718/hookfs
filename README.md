# HookFS: A Usermode Hookable Filesystem Library

[![GoDoc](https://godoc.org/github.com/osrg/hookfs?status.svg)](https://godoc.org/github.com/osrg/hookfs)

## 丰富入参

    Usage of ./faultfs:
    ./faultfs [OPTIONS] MOUNTPOINT ORIGINAL...
    Options
      -delay duration
            delay time (ms, s, min, h ...) ,if use type (ReadFileDelay , WriteFileDelay, FsycnDelay), must set this option
      -log-level int
            log level (0..2)
      -mountpoint string
            (required) mount point , must set this option (default "/mnt/faultfs")
      -original string
            (required) original dir, must set this option (default "/mnt/fs")
      -percent int
            fault percentage (0, 99]
      -type int
            (required)type value, type list:
               0   OpenFileEIO
               1   OpenFileEPERM
               2   ReadFileDelay
               3   ReadFileErr
               4   WriteFileENOSPC
               5   WriteFileDelay
               6   MkDirEACCES
               7   MkDirEPERM
               8   RmDirEACCESS
               9   RmDirEPERM
               10  FsycnDelay
               11  FsycnEIO
               12  OpenDirEACCESS
               13  OpenDirEPERM    

* 示例

```
./faultfs -mountpoint /root/testfs/test-dir-1 -original /root/testfs/test-dir-2 -type 13 -percent 90 -delay=2s
```

## Possible Recipes

* Fault Injection (Example: [Namazu](https://github.com/osrg/namazu))
* Cache
* Malware Detection

and so on..

HookFS was originally developed for [Namazu](https://github.com/osrg/namazu), but we believe HookFS can be also used for other purposes.

## Install

    $ go get github.com/osrg/hookfs/hookfs

## Running an Example

    $ cd example/ex01
    $ go build
    $ ./ex01 "/mnt/hookfs" "/original"
	^C
    $ fusermount -u "/mnt/hookfs"

## API Design
You have to implement `HookXXX` (e.g. `HookOnOpen`, `HookOnRead`, `HookOnWrite`, ..)  interfaces.

```go
type HookOnRead interface {
	// if hooked is true, the real read() would not be called	
	PreRead(path string, length int64, offset int64) (buf []byte, hooked bool, ctx HookContext, err error)
	PostRead(realRetCode int32, realBuf []byte, prehookCtx HookContext) (buf []byte, hooked bool, err error)
}
```
	
Then, regist your hook implementation to the HookFS server.

```go
fs, err := NewHookFs("/original", "/mnt/hookfs", &YourHook{})
if err != nil { .. }
err = fs.Serve()
```

See [`hook.go`](src/hookfs/hook.go) for further information. [GoDoc](https://godoc.org/github.com/osrg/hookfs) is also your friend.

## Related Projects
* [Namazu (Earthquake)](https://github.com/osrg/namazu)
* [HookSwitch](https://github.com/osrg/hookswitch)

## How to Contribute
We welcome your contribution to HookFS.
Please feel free to send your pull requests on github!

## Copyright
Copyright (C) 2015 [Nippon Telegraph and Telephone Corporation](http://www.ntt.co.jp/index_e.html).

Released under [Apache License 2.0](LICENSE).
