# doke

run make tasks in docker. s/make/doke/

## When to use doke?

If you have a Makefile, and you want to make it portable and repeatable - use `doke` instead of `make`.

This can be useful for synching up environments across a project team.

Note: if you don't already use `make`, there's probably better options out there for you. Note that doke does nothing to isolate individual containers for different build steps (I'd need to write a Makefile parser for that) - try `docker-compose`, `docker`, or something like that.

## How?

 * Install doke
 
 ```
    go get -u -v github.com/laher/doke
 ```

 * Create a Dockerfile.doke

This example is for a go app ('myapp'), which uses `jq` as part of the build process.

```
  FROM golang:alpine 

  RUN mkdir /myapp
  RUN apk add --update make jq

  WORKDIR /myapp

  ENTRYPOINT ["make"]
```

 * Now run `doke` instead of `make`

```
   doke build
```
   
## And then?

Well, now you can rest assured that your build is portable - any unusual dependencies can now be encapsulated in your Dockerfile.doke

## Next steps

I'll make a bunch of example Dockerfiles to help you get started with common tech stacks.
