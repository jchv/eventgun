# EventGun
EventGun is a tool for sending filesystem events over the network. Prime use cases would include networked filesystems, such as when using Docker, where inotify does not work properly.

## Concept
On the client-side, EventGun contains a library that intercepts API calls; currently only Linux and inotify are supported. This is done by using hooking mechanisms. On Linux, we produce a shared library (`eventgun.so`) that can be used with `LD_PRELOAD`. Once this library is injected into a process, attempts to call the `inotify_init` syscall will trigger EventGun to connect to the event server.

On the server-side, the EventGun server will watch for filesystem changes locally and translate those events such that they can be interpreted by the client, by translating the paths as necessary.
