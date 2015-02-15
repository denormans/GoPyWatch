GoPyWatch
=========

Overview
--------

GoPyWatch is a utility that executes a Python file and watches it for changes.

    Usage:
      gopywatch --interactive --file someapp.py --extraWatchDir somedir


OS Support
----------

It should pick up file changes automatically for OS X. Feel free to contribute support for other OSes if you find one that doesn't work.

Dependencies
------------

    go get denormans/gopywatch

Build
-----

    go build -o bin/gopywatch denormans/gopywatch/main

Run
---

    go run denormans/gopywatch/main/main.go --interactive --file someapp.py --extraWatchDir somedir

Feedback
--------

If you'd like to contribute, just submit a pull request. If you find any issues, feel free to add an issue, but it'll get fixed faster if you fix it yourself :)

Background
----------

I wrote this app to help my son learn Python.

We are using the [Turtle Graphics] module, which requires interactive mode or some sort of pause at the end of the program
so the graphical window doesn't disappear immediately after drawing. This can be annoying when wanting to restart the program
with every change.

I was surprised to find that there was no Python watcher like there is for [NodeJS] eg [nodemon]. Therefore, I decided to learn [Go Lang] and build one.

I totally understand this initial implementation is naive and minimal.


[Go Lang]: http://golang.org/
[Python]: https://www.python.org/
[Turtle Graphics]: https://docs.python.org/2/library/turtle.html
[NodeJS]: http://nodejs.org/
[nodemon]: http://nodemon.io/

