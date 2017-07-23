[![Build Status](https://secure.travis-ci.org/erriapo/stats.png)](https://travis-ci.org/erriapo/stats)

[![GoDoc](https://godoc.org/github.com/erriapo/stats?status.png)](https://godoc.org/github.com/erriapo/stats)

An implementation of B.P. Welford's algorithm to maintain a running variance
over a stream of observations. 

This code is based on a John D Cook's [blogpost](https://www.johndcook.com/blog/standard_deviation/).

## Example usage

```go
s := stats.NewSink() 
s.Push(-1.1813782)
s.Push(-0.2449577)
s.Push(0.8799429)

fmt.Printf("Mean = %v\n", s.Mean()) // -0.182131...
fmt.Printf("Standard Deviation = %v\n", s.StandardDeviation()) // 1.032095...

```

## Todos

* Add a method that returns the running median.
* Make the implementation multi-goroutine safe.

## License

* Code is released under MIT license. 
