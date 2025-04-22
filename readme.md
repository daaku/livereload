# livereload

Package livereload provides some JavaScript and corresponding `http.Handler`
to enable full page reloads when the `http.Server` restarts.

## demo

Pairing with [gow](https://github.com/mitranim/gow) is probably the easiest
way to use this. Look at the [demo](demo/demo.go) code and run it like such:

```sh
gow run ./demo
```

Then make a change to the HTML [demo](demo/demo.go) in the `index` function
and see your browser automatically reload the page.
