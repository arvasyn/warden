#!/usr/bin/env python3
import os
import subprocess

if os.path.exists("/home/marina"):
    print("I have access")
else:
    print("I don't have access")

def run_shell():
    print("Shell")
    while True:
        try:
            cmd = input(f"{os.getcwd()} $ ").strip()

            if not cmd:
                continue
            if cmd.lower() in ['exit', 'quit']:
                break

            if cmd.startswith("cd "):
                path = cmd.split(" ", 1)[1]
                os.chdir(os.path.expanduser(path))
                continue

            subprocess.run(cmd, shell=True)

        except KeyboardInterrupt:
            print("\nUse 'exit' to quit.")
        except Exception as e:
            print(f"Error: {e}")

run_shell()
