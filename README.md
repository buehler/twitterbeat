# Twitterbeat

Welcome to Twitterbeat.

Ensure that this folder is at the following location:
`${GOPATH}/github.com/buehler`

## Getting Started with Twitterbeat

### Init Project
To get running with Twitterbeat, run the following commands:

```
glide update --no-recursive
make update
```


To push Twitterbeat in the git repository, run the following commands:

```
git init
git add .
git commit
git remote set-url origin https://github.com/buehler/twitterbeat
git push origin master
```

For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).

### Build

To build the binary for Twitterbeat run the command below. This will generate a binary
in the same directory with the name twitterbeat.

```
make
```


### Run

To run Twitterbeat with debugging output enabled, run:

```
./twitterbeat -c twitterbeat.yml -e -d "*"
```


### Test

To test Twitterbeat, run the following commands:

```
make testsuite
```

alternatively:
```
make unit-tests
make system-tests
make integration-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`


### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `etc/fields.yml`.
To generate etc/twitterbeat.template.json and etc/twitterbeat.asciidoc

```
make update
```


### Cleanup

To clean  Twitterbeat source code, run the following commands:

```
make fmt
make simplify
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


### Clone

To clone Twitterbeat from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/github.com/buehler
cd ${GOPATH}/github.com/buehler
git clone https://github.com/buehler/twitterbeat
```


For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).
