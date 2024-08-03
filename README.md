# Unbounded Golang Channels

`unbounded` is a Go package that provides an implementation of an unbounded/unlimited/unsized channel with chunked buffers and a wait list. This implementation allows for efficient inter-goroutine communication with unbounded capacity, ensuring smooth performance even under high concurrency without blocking on send for a receiver.

## Features

- Unbounded capacity: The channel grows dynamically as needed.
- Efficient memory usage: Uses chunked buffers to manage items.
- Wait list: Handles blocking receives without busy-waiting.
- No blocking send: Send will not block sender for a receiver.

## Installation

To install the package, run:

```bash
go get github.com/fereidani/unbounded
```

## Usage

### Creating an Unbounded Channel

To create a new unbounded channel, use the `New` function:

```go
import "github.com/fereidani/unbounded"

ch := unbounded.New[int]()
```

### Sending Items

To send an item to the unbounded channel, use the `Send` method:

```go
ch.Send(777)
```

### Receiving Items

To receive an item from the unbounded channel, use the `Receive` method:

```go
value, ok := ch.Receive()
if ok {
    fmt.Println("Received:", value)
} else {
    fmt.Println("Channel closed")
}
```

### Closing the Channel

To close the unbounded channel, use the `Close` method:

```go
ch.Close()
```

## Example

Here is a complete example demonstrating how to use the `unbounded` package:

```go
package main

import (
    "fmt"
    "github.com/fereidani/unbounded"
)

func main() {
    ch := unbounded.New[int]()

    go func() {
        for i := 0; i < 10; i++ {
            ch.Send(i)
        }
        ch.Close()
    }()

    for {
        value, ok := ch.Receive()
        if !ok {
            break
        }
        fmt.Println("Received:", value)
    }
}
```

## Benchmarks

Benchmarks have been conducted to compare the performance of the unbounded channel implementation with Go's built-in channels.
For some reason this library is outperforming Golang channels in SPSC.

```
BenchmarkUnboundedChannelMPMC-8   	7940866	       166.4 ns/op
BenchmarkGoChannelMPMC-8         	10961176	       105.3 ns/op
BenchmarkUnboundedChannelSPSC-8  	21728120	        56.32 ns/op
BenchmarkGoChannelSPSC-8         	15320703	        66.16 ns/op
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
