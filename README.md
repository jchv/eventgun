# EventGun
EventGun is a tool for sending filesystem events over the network. Prime use cases would include networked filesystems, such as when using Docker, where inotify does not work properly.

## Status
This project is not functioning as intended yet, please check back later. Currently it seems something is wrong with the file descriptor it sends to other applications, they freeze up when trying to read an entry. Also, the code is incredibly hasty.

## Concept
On the client-side, EventGun contains a library that intercepts API calls; currently only Linux and inotify are supported. This is done by using hooking mechanisms. On Linux, we produce a shared library (`eventgun.so`) that can be used with `LD_PRELOAD`. Once this library is injected into a process, attempts to call the `inotify_init` syscall will trigger EventGun to connect to the event server.

On the server-side, the EventGun server will watch for filesystem changes locally and translate those events such that they can be interpreted by the client, by translating the paths as necessary.