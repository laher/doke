# doke

Run your make tasks in docker with a drop-in command. `s/make/doke/`

## When to use doke?

If you have a Makefile, and you want to make it portable and repeatable - use `doke` instead of `make`.

This can be useful for mitigating environment differences across a project team. The "it works on my laptop" problem.

_Note: if you don't already use `make`, there's probably better options out there for you. Note that doke does nothing to isolate individual containers for different build steps (for this I'd need to really replace make in your container) - try `docker-compose`, `drone`, or something like that._

## How?

 * Install doke
 
 ```
    go get -u -v github.com/laher/doke
 ```

 * Run `doke` once to generate a special dockerfile for doke to build with:

```
   doke
```
   
This creates a file called Dockerfile.doke, which doke will use to execute your make tasks. It's super minimal, your make tasks will fail (for now).

 * Edit `Dockerfile.doke`

If you're happy to use `alpine`, you can get a long way with more dependencies - `RUN apk add ...`.

Otherwise, you can start FROM any image (e.g. `ubuntu`), as long as you remember to install `make` and any other build dependencies.

_Notice how doke's own [Dockerfile.doke](Dockerfile.doke) uses golang:alpine instead of straight alpine._

 * Now run `doke` instead of `make`

```
   doke build
```

## Env vars

Env vars in docker are passed in explicitly. This is for the best, but it means doke works a little bit differently from make:

Instead of:

```
MY_ENV_VAR=1 make print-an-env-var
```

Use the `-e` flag instead:

```
doke -e MY_ENV_VAR=1 print-an-env-var
```

## And then?

Well, now you can rest assured that your build is portable - any unusual dependencies can now be encapsulated in your Dockerfile.doke

_NOTE: `doke` doesn't currently let you actually install something outside your working directory._

## Next steps

 * I'll make a bunch of example Dockerfiles to help you get started with common tech stacks.
 * I might add a 'dokefile' config file, to allow for mapping extra volumes for certain tasks (e.g. `make install`). I'm not sure yet about security considerations so maybe not.

## Probably not

 * I've considered some fancy stuff to parse the Makefile using [mmake's parser](https://godoc.org/github.com/tj/mmake/parser) to isolate steps use separate containers for each task.
It sounds messy and complicated though, so I probably won't.
