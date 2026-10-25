# view-nux

A tiny `top`-like process viewer, written twice — once in Go, once in Python —
as a learning exercise in how Linux exposes process info and how a core OS
concept (process state, memory, command line) gets read straight from the
kernel without any external libraries.

## What it does

Every second it clears the screen and prints a table of every running
process:

```
    PID STATE   RAM(MB)  COMMAND
      1 S           2.3  /sbin/init
    412 S           8.1  /usr/bin/dockerd
   1337 R           0.4  python3 main.py
```

Press `Ctrl+C` to exit.

## How it works

Linux exposes a live view of every process as files under `/proc/<pid>/`.
There's no syscall needed — it's just text files the kernel keeps in sync:

| File | What we pull from it |
|---|---|
| `/proc/<pid>/stat` | process name (in parentheses) and state (`R` running, `S` sleeping, `D` disk wait, `Z` zombie, etc. — it's the 3rd field, but the name itself can contain spaces/parens, so we locate it by the last `)` rather than splitting naively) |
| `/proc/<pid>/statm` | memory usage in **pages**; multiply by the page size to get bytes |
| `/proc/<pid>/cmdline` | the full command line, NUL-separated instead of space-separated |

The directory listing of `/proc` itself is the process list: every entry
that's purely numeric is a PID, everything else (`net`, `self`, `cpuinfo`,
...) is skipped.

Page size matters because `statm` reports in pages, not bytes:
- Go: `os.Getpagesize()`
- Python: `os.sysconf("SC_PAGE_SIZE")`

Both are almost always 4096 bytes on x86/ARM Linux, but reading it rather
than hardcoding it is the correct habit.

## Files

- `main.go` — Go version
- `main.py` — Python version

Both do the same thing; comparing them side by side is a good way to see
where Go makes you handle errors and types explicitly that Python lets you
wave away (until an exception does it for you at runtime instead).

## Running it

**Go:**
```bash
go run main.go
```
or build a binary:
```bash
go build -o taskviewer main.go
./taskviewer
```

**Python:**
```bash
python3 main.py
```

No dependencies beyond the standard library on either side.

## Ideas to extend it

Things worth trying next, roughly in order of difficulty:

- Sort output by RAM or PID (`-cpu`/`-mem` style flags)
- Add CPU% by sampling `/proc/<pid>/stat` utime/stime twice and diffing
- Filter by process name or a minimum RAM threshold
- Add color (green for `R`, yellow for `D`, red for `Z` zombies)
- Replace the polling loop with a TUI library (`tview` for Go,
  `curses`/`rich` for Python) instead of `clear`-and-reprint
- Handle the race condition where a process exits between `os.ReadDir`
  and reading its files (currently just skipped via the `err != nil`
  checks — that's actually the correct fix, worth noticing)

## Notes / gotchas

- Needs to run on Linux (or WSL) — `/proc` doesn't exist on macOS or
  native Windows.
- You won't see full details for processes owned by other users unless
  you run as root — this is the kernel enforcing permissions, not a bug
  in the code.
- Zombie processes (`Z`) will show up with no meaningful command line;
  that's expected, they're already dead and just waiting to be reaped.
