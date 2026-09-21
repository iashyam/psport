# psport

A small CLI that tells you what's listening on a given port. Wraps `lsof` and prints a clean table instead of raw output.

## Install

### Using `go install` (recommended)

```bash
go install github.com/iashyam/psport@latest
```

This downloads, builds, and drops the `psport` binary into `$GOPATH/bin` (usually `$HOME/go/bin`). Make sure that directory is on your `$PATH`:

```bash
# bash/zsh — add to ~/.bashrc or ~/.zshrc
export PATH="$PATH:$(go env GOPATH)/bin"
```

Verify it's installed:

```bash
psport --help 2>/dev/null; which psport
```

### Build from source

```bash
git clone git@github.com:iashyam/psport.git
cd psport
go build -o psport .
sudo mv psport /usr/local/bin/   # or any dir on your $PATH
```

## Usage

```bash
psport [-q|--quiet] [-f|--free] <port>
```

Example:

```bash
$ psport 3000
COMMAND  PID   USER   PROTO  ADDRESS          STATE
node     1234  shyam  TCP    *:3000           LISTEN
```

If nothing is listening on the port, it prints a message and exits cleanly.

### Quiet mode

`-q`/`--quiet` prints just the PID(s), one per line, no header — handy for piping into `kill`:

```bash
$ psport -q 3000
1234

$ kill $(psport -q 3000)
```

On an empty port, quiet mode prints nothing and exits cleanly.

### Free mode

`-f`/`--free` prints the table, then asks a single confirmation before killing every process listening on the port:

```bash
$ psport -f 3000
COMMAND  PID   USER   PROTO  ADDRESS          STATE
node     1234  shyam  TCP    *:3000           LISTEN

Kill 1 process and free the port? [y/N]: y
killed PID 1234 (node)
```

Answering anything other than `y`/`yes` aborts and kills nothing.

## Requirements

- `lsof` must be installed and on your `$PATH` (present by default on macOS and most Linux distros).
- Go 1.27+ to build from source.

## How it works

`psport` shells out to `lsof -nP -i :<port> -F pcuLPnT`, parses its machine-readable field output, and renders it as an aligned table.

---

Source at [github.com/iashyam/psport](https://github.com/iashyam/psport) · made by [iashyam](https://github.com/iashyam)
