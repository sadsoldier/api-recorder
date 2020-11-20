## dummy-api

Simple data logging proxy

### Configuration

dummyapi.yml:

- port: [integer 1-65536]
- target: [string http://some.hostname.org:port]

Example:
```
port: 8017
target: http://www.gnu.org:80
```

### Limitation

- Only HTTP
- No for stream (big data request)

### Installation and start

```
$ git clone ...
$ cd dummy-api
$ ./configure --prefix=/usr/local
$ make 
$ make install
```

For run inplace
```
$ ./configure --enable-devel-mode
$ make
$ ./dummyapi 
```

### Options 

```
$ ./dummyapi --help

usage: dummyapi command [option]

  -d    debug mode
  -debug
        debug mode
  -devel
        devel mode
  -e    devel mode
  -f    run in foreground
  -foreground
        run in foreground
  -p int
        listen port (default 8050)
  -port int
        listen port (default 8050)
```

