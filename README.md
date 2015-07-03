Lemonade
========

remote...lemote...lemode......Lemonade!!! :lemon: :lemon:

Lemonade is a remote utility tool.
(copy|paste|open browser) over TCP.


Installation
------------

```sh
go get github.com/pocke/lemonade
```

~~Or download from latest release~~ (TODO)


Usage
-----

Default port is `2489`.


### Server

```sh
lemonade server --port 1234
```

### Client

`192.168.xx.xx` is a IP address of lemonade Server.

#### Open browser

```sh
lemonade open --port 1234 --host '192.168.xx.xx' 'http://example.com'
# or
echo 'http://example.com' | lemonade open --port 1234 --host '192.168.xx.xx' 
```


`http://example.com` is opened by browser on Server.


#### Copy

```sh
lemonade open --port 1234 --host '192.168.xx.xx' 'hogefuga'
# or
echo 'hogefuga' | lemonade open --port 1234 --host '192.168.xx.xx'
```

`hogefuga` is copied to Server.


#### Paste

```sh
lemonade open --port 1234 --host '192.168.xx.xx'
# => hogefuga
```

`hogefuga` is a clipboard value of Server.



Configuration
--------------

You can override command line options by configuration file.

### Server

`~/.config/lemonade.toml`

```toml
port = 1234
allow = '192.168.0.0/24'
```

`allow` is a comma separated list of allowed IP(and subnet mask) or an allowed IP.  
Default value is `0.0.0.0/0,::0`(allowed from all IP).

### Client

`~/.config/lemonade.toml`

```toml
port = 1234
host = '192.168.x.x'
```


Alias
-----

You can use lemonade as a `xdg-open` , `pbcopy` and `pbpaste`.

For example.

```sh
ln -s /path/to/lemonade /usr/bin/xdg-open
xdg-open  'http://example.com' # Same as lemonade open 'http://example.com'
```
