### What

This is a simple helper tool which lets you `hijack` using a [Concourse](http://concourse.ci) URL. Instead of doing something like:
```
fly -t <target> hijack -j <pipeline>/<job>
```

This tool will let you use a Concourse URL to hijack the relevant container. For example:

```
fly-hijack http://concourse.acme.org:8080/teams/main/pipeline/production/jobs/BuildAll/builds/27
```

This way lets you simply grab the Concourse URL and paste it to the `fly-hijack` command.

The tool assumes you're using Concourse >= 2.x which introduced teams. Also assumes your `fly` binary is at `/usr/local/bin/fly`.

It uses the given URL to infer the correct target by matching against what you have in `~/.flyrc`. You can also provide an explicit target using the `-t` option.

### Building

Requires using [govendor](https://github.com/kardianos/govendor) for dependencies.

```
go get github.com/kardianos/govendor
```

Assuming you've cloned this repo, then do:

```
cd src/fly-hijack
govendor sync
cd -
go build fly-hijack
```