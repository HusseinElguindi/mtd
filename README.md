# mtd
Go multithreaded/concurrent downloader.

## Installation
```bash
go get github.com/husseinelguindi/mtd
```

## Features
- [x] Concurrent chunked download through Goroutines
- [x] Offset write to file
- [x] Synchronized writing to file (never write 2 things at once)
- [x] Go module interface to download within your program
- [ ] Command line interface
    - [ ] Args
    - [ ] Console input

## Design Choices
### Synchronized Writing
- The idea was that moving (seeking in code) the disk head of a hard drive, slows the writing of a file.
- Writing larger chunks at a time would yield better performance, avoiding the hard drive seek delay.
