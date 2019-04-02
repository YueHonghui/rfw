rfw(simple rotate file writer)
===
## What is rfw
Package rfw is a simple rotate file writer that can rotate file dayly and be configured how many rotated files to remain. It implements `io.WriteCloser` and can be used with many logger package easyly.

## Why to use rfw
There are many logger package in golang ecosystem, but few provide ability to rotate log file. Though we can use logrotate to rotate log file, but it has some drawbacks
1. when logrotate working, load of io increasing too(few of programs implements a custom signal to reopen log file)
1. risk of lost some lines of log content
1. difficult to use

## How to use
See code sample below
```golang
    // file name of rotated log file will be ./logfile-20190101 ...
    logpath := "./logfile"
    remaindays := 7
	lw, err = rfw.NewWithOptions(logpath, rfw.WithCleanUp(remaindays))
	if err != nil {
        //...
    }
    
    //now you have a io.WriteCloser lw, you can use it with std logger
    log.SetOutput(lw)
```

## Status of this package
This package has been used in production for more than two years. It is ready for production use.

## Roadmap
1. [TODO] async log writing