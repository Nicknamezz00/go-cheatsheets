## Test
```shell
go test unsafemap_test.go --race
```

Failed, race detected.

```shell
go test safemap_test.go safemap.go --race
```

Passed.
