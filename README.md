### buildkit-pack

BuildKit frontend for building buildpacks directly.


#### Usage

With Docker (v18.06+ with `DOCKER_BUILDKIT=1`):

Add `# syntax = tonistiigi/pack` as the first line of a file (eg. `manifest.yml`):

```
docker build -f manifest.yml .
```

With `buildctl`:
```
buildctl build --frontend=gateway.v0 --frontend-opt source=tonistiigi/pack --local context=.
```
