## srv

Srv is a simple static file server. It accepts a directory path as input
from which it will serve any requested file.

It optionally supports TLS encryption using Let's Encrypt and autocert.
The server supports clean shutdown by trapping OS signals and subsequently
forcing the serve to shut down gracefully.

The server comes with sane read/write and idle timeouts to ensure safe
behaviour when exposed to the internet.

When run on Linux, the server supports the use of systemd sockets by specifying
the listener address as a socket name, prefixed with `"systemd:"`.


## Usage

    $ go get github.com/hexaflex/srv


To serve contents in the current working directory on localhost:

    $ srv

To serve a specific directory on a fixed address/port:

    $ srv -addr ":8080" assets/

To use TLS encryption, specify the `-tls` flag. Note that this will only work
when the server is accessible through a public domain:

    $ srv -tls -addr "mydomain.org:8080" assets/

When using TLS, certificates are cached in `$TEMP/srv-certs`.


## License

Unless otherwise stated, this project and its contents are provided under a 3-Clause BSD license.
Refer to the LICENSE file for its contents.