# mtd
Golang multithreaded/concurrent downloader.
Work in progress, any contribution or feedback is appreciated!

## Installation
```bash
go get github.com/husseinelguindi/mtd
```

## Features/To-Do
- [x] Concurrent chunked download through Goroutines
- [x] HTTP/HTTPS download support
- [x] Offset write to file
- [x] Synchronized writing to file (never write 2 things at once)
- [x] Go module/library interface
- [x] Multiple downloads in a single instance
- [ ] Pause/resume support
- [ ] Command line interface
    - [ ] Args
    - [ ] Console input

## Design Choices
### Synchronized Writing
A writer Goroutine that only writes one chunk at a time, this was done for a number of a simple reason:
- Seeking to a location in a file, physically moves the write head of a hard drive, which slows the writing to a file.
- This means that writing chunks, that are near each other, would yield better performance, avoiding the hard drive seek delay.
- This also means that using a bigger chunk/buffer size is beneficial to write speeds (less system calls, as well).