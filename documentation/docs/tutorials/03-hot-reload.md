# Hot reload

Hot reload is a feature that allows you to update your code and
see the changes in real-time without restarting the server.
This is very useful for development,
as it allows you to see the changes you make to your code immediately.

To enable hot reload, you need to install the `air` command-line tool:

```sh
go install github.com/cosmtrek/air@latest
```

Then, create a `.air.toml` file in the root of your project with the following content:

```sh
air init
```

Finally, simply the following command to start the server with hot reload:

```sh
air
```
