# doke

Run your make tasks via docker (or docker-compose) with a drop-in command. `s/make/doke/`

doke can really be seen as a starter pack for migrating your build pipeline to use containers. It's pretty dumb. 

## When to use doke?

If you have a Makefile, and you want to make it portable and repeatable - use `doke` instead of `make`.

This can be useful for mitigating environment differences across a project team. The "it works on my laptop" problem.

_Note: if you don't already use `make`, there's probably better options out there for you. Note that doke does nothing to isolate individual containers for different build steps (for this I'd need to really replace make in your container) - try `docker-compose`, `drone`, or something like that._

## How?

 * Pre-requisites - docker and docker-compose
   * Install [docker](https://www.docker.com/get-started)
   * [Optional]: install [docker-compose](https://docs.docker.com/compose/install/)
   (docker-compose is optional, but recommended)

### Install doke
 
 ```
    go get -u -v github.com/laher/doke
 ```

### Setup

 * Run `doke` once to generate a special dockerfile for doke to build with:

```
   doke
```
   
This creates some files, which doke will use to execute your make tasks in docker[-compose]. It's super minimal, your make tasks will almost definitely fail (for now). The filenames are named `*.doke*` to avoid naming collisions. You can go ahead and rename them and use the standard tools.

### Update configs

 * Edit `Dockerfile.doke` and/or `docker-compose.doke.yaml`

#### Dockerfile: 

 * If you're happy to use `alpine`, you can get a long way with more dependencies - `RUN apk add ...`. 
 * For many technologies you may like to start with a recommended image (e.g. `node:11-stretch`). _Notice how doke's own [Dockerfile.doke](Dockerfile.doke) uses `golang:alpine` instead of straight `alpine`._
 * Otherwise, you can start FROM any image (e.g. `ubuntu`), as long as you remember to install `make` and any other build dependencies.

#### docker-compose file: 

 * If your task is failing because it can't access a file outside your working dir, you can update the volume to use a parent dir (e.g. `$PWD/../..`), and then update the `working_dir` accordingly
 * Add other containers to provide services like databases


 * Now run `doke` instead of `make`

```
   doke build
```

#### Env vars

Env vars in docker are passed in explicitly. This is for the best, but it means doke works a little bit differently from make:

Instead of:

```
MY_ENV_VAR=1 make print-an-env-var
```

Use the `-e` flag instead:

```
doke -e MY_ENV_VAR=1 print-an-env-var
```

Alternatively you can specify them in your docker-compose file

## And then?

Well, now you can rest assured that your build is portable - any unusual dependencies can now be encapsulated in your Dockerfile.doke

_NOTE: `doke` doesn't currently let you actually install something outside your working directory._

## Next steps

 * I'll make a bunch of example Dockerfiles to help you get started with common tech stacks.
 * `--env-file` support
 * docker-compose: support for mounting your parent dir as the main volume (especially when your Makefile isn't in the root of your repo)

## Probably not (out of scope)

 * I've considered some fancy stuff to parse the Makefile using [mmake's parser](https://godoc.org/github.com/tj/mmake/parser) to isolate steps use separate containers for each task.
It sounds messy and complicated though, so I probably won't.
