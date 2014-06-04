scripts2http
====

Inspired by [pyjojo](https://github.com/atarola/pyjojo), script2http exposes a directory of scripts as http get endpoints with the smallest possible functionality! (because I needed something that could work without dreaded python installations)

There are many used, but this was build to have a simple docker container instrocpection endpoint, e.g: inside a container you could run `curl -s 172.17.42.1/container-data/$(hostname)` and return the result of `docker inspect $1` where `$1` would equal the passed hostanme in the URI.

To try it locally in development: `go run script2http.go -scripts-dir sample-scripts`

Thanks @atarola for the wonderful idea.