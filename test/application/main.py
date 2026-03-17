#!/usr/bin/env python3
import os
import tkinter as tk

root = tk.Tk()
root.title("Hello World")
tk.Label(root, text="Hello, World!").pack(padx=20, pady=20)


if os.path.exists("/home/marina"):
    print("I have access")
else:
    print("I don't have access")

root.mainloop()
