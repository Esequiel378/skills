#!/usr/bin/env python3
"""Print Codex's live rate limits as JSON: weekly used percent and reset time."""
import json
import subprocess
import sys
import threading

p = subprocess.Popen(["codex", "app-server"], stdin=subprocess.PIPE, stdout=subprocess.PIPE,
                     stderr=subprocess.DEVNULL, text=True)
threading.Timer(20, p.kill).start()


def send(m):
    p.stdin.write(json.dumps(m) + "\n")
    p.stdin.flush()


send({"id": 1, "method": "initialize", "params": {"clientInfo": {"name": "batch-orchestrate", "version": "0"}}})
send({"method": "initialized"})
send({"id": 2, "method": "account/rateLimits/read"})
for line in p.stdout:
    m = json.loads(line)
    if m.get("id") == 2:
        print(json.dumps(m.get("result") or m.get("error")))
        p.kill()
        sys.exit(0 if "result" in m else 1)
sys.exit(1)
