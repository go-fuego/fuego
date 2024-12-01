# Hot reload

Hot reload is a feature that allows you to update your code and
see the changes in real-time without restarting the server.
This is very useful for development,
as it allows you to see the changes you make to your code immediately.

To enable hot reload, you need to install the `air` command-line tool.

```sh
go install github.com/air-verse/air@latest
```

Optionally, create a `.air.toml` configuration file to customize the hot reload behavior.

```sh
air init
```

Simply the following command to start the server with hot reload.

```sh
air
```
