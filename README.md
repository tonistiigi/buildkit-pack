[![asciicinema example](https://asciinema.org/a/2huuQ3iYTAXTGrrLRD4xd50XS.png)](https://asciinema.org/a/2huuQ3iYTAXTGrrLRD4xd50XS)


## buildkit-pack

BuildKit frontend for building buildpacks directly.


### Usage

#### With Docker (v18.06+ with `DOCKER_BUILDKIT=1`):

Add `# syntax = tonistiigi/pack` as the first line of a file (eg. `manifest.yml`):

```
docker build -f manifest.yml .
```

#### With `buildctl`:
```
buildctl build --frontend=gateway.v0 --frontend-opt source=tonistiigi/pack --local context=.
```

### Options

Detection can be enabled with `skipDetect=true`. Custom buildpacks can be set in `manifest.yml` or passed with `buildpackOrder=url`.


### Examples


#### Python

```
> git clone git://github.com/ihuston/python-cf-examples.git
> cd python-cf-examples/01-simple-python-app
> tmp=$(mktemp); ( echo "# syntax = tonistiigi/pack" ; cat manifest.yml ) > $tmp && mv $tmp manifest.yml
> docker build -t pythonapp -f manifest.yml .
[+] Building 89.4s (16/16) FINISHED
 ....
 => => writing image sha256:1e13eb221e7779c8aa65624b97afcc11ee797d3c09d80dcfc1e55d956d467d39         0.0s
 => => naming to docker.io/library/pythonapp                                                         0.0s
> docker run -d -p 8080:8080 pythonapp
> curl localhost:8080
Hello World! I am instance 0
```

#### Go

```
> git clone git://github.com/cloudfoundry/go-buildpack
> cd go-buildpack/fixtures/go_app/
> tmp=$(mktemp); ( echo "# syntax = tonistiigi/pack" ; cat Procfile ) > $tmp && mv $tmp Procfile
> docker build -t goapp -f Procfile .
[+] Building 28.9s (16/16) FINISHED
...
 => => writing image sha256:f0d240bed9a2b12b288845c3a305050c8d9fa035f564d057e95b0714f68363f2         0.0s
 => => naming to docker.io/library/goapp                                                             0.0s
> docker run -d -p 8080:8080 goapp
> curl localhost:8080
go, world
.

.
```