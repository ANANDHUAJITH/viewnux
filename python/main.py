
import os
import time

def tasks():
    print(f"{'PID':>7} {'STATE':>5} {'RAM(MB)':>9}  COMMAND")

    for pid in os.listdir("/proc"):
        if not pid.isdigit():
            continue

        try:
            with open(f"/proc/{pid}/stat") as f:
                stat = f.read()

            name = stat[stat.find("(") + 1:stat.rfind(")")]
            fields = stat[stat.rfind(")") + 2:].split()
            state = fields[0]

            with open(f"/proc/{pid}/statm") as f:
                ram = int(f.read().split()[1])
                ram = ram * os.sysconf("SC_PAGE_SIZE") / 1048576

            with open(f"/proc/{pid}/cmdline", "rb") as f:
                cmd = f.read().replace(b"\0", b" ").decode(
                    errors="replace"
                ).strip()

            if not cmd:
                cmd = f"[{name}]"

            print(f"{int(pid):7} {state:>5} {ram:9.1f}  {cmd}")

        except (OSError, ValueError, IndexError):
            continue


try:
    while True:
        os.system("clear")
        tasks()
        time.sleep(1)
except KeyboardInterrupt:
    print("\nExited.")
