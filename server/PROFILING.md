## Guide to CookieServer profiling


1. Build the CookieServer binary with profiling flags:

`@go build -gcflags='github.com/ByteTheCookies/cookieserver/...="-m"' -o $(BIN_DIR)/$(BINARY_NAME) $(CMD_DIR)/$(MAIN_FILE)`

2. Run the CookieServer binary with profiling enabled:
3. Start the Pyroscope server:
```sh
docker pull grafana/pyroscope:latest
docker run -d -p 4040:4040 grafana/pyroscope:latest
```

4. Start the Pyroscope agent:

`go get github.com/grafana/pyroscope-go`

```go

// e poi nel main in debug
pyroscope.Start(pyroscope.Config{
		ApplicationName: "simple.golang.app",
		ServerAddress:   "http://pyroscope-server:4040",
		Logger:          pyroscope.StandardLogger,

		ProfileTypes: []pyroscope.ProfileType{
			// these profile types are enabled by default:
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,

			// these profile types are optional:
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
	})
```
