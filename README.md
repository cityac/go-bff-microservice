# BFF


## Getting started

`Create env`
```
cp .env.example .env
```

`Install demon`

- [ ] [Install gow](https://github.com/mitranim/gow)

`Enable download of sms-core library`

```
go env -w GOPRIVATE=gitlab.HIDDEN.com/HIDDEN/sms-core

git config --global url."git@gitlab.HIDDEN.com:".insteadOf "https://gitlab.HIDDEN.com/"
```

`Fetch submodules`
```
git submodule update --init
```

`Run project`
```
gow run cmd/main.go
```

# Profiling

To view profiling board

```
http://localhost:8080/debug/pprof/
```
To view running goroutines
```
http://localhost:8080/debug/pprof/goroutine?debug=1
```