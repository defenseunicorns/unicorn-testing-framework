# dos-games Package Testing Example

This is an example zarf package with custom package testing built into the package.

NOTE: manifests and parts of zarf.yaml copied from [zarf#examples/dos-games](https://github.com/defenseunicorns/zarf/tree/main/examples/dos-games)

## Quick Start

From this directory
1. `go test -c -o ./dos-games-test` #Build the Test Code
2. `zarf package create` #Build the Zarf Package
3. `zarf init && zarf package deploy` # deploy the package

On successful completion you should see similar output:

```sh
 ✔  Processing helm chart raw-dos-games-baseline-multi-games:0.1.1676505567 from Zarf-generated helm chart
     === RUN   TestGamesAreRunning
         package_test.go:14: Running test: TestGamesAreRunning
         package_test.go:40: The address we received is: http://127.0.0.1:55826
     --- PASS: TestDoomIsRunning (2.34s)
     PASS
  ✔  Completed command "./dos-games-test -test.v"
```

## Explanation

`go test -c ...`
    Test code can be compiled into a binary!

`zarf.yaml`
Adding these lines to the zarf package includes the test binary and defines the zarf onDeploy actions
```
    files:
      - source: dos-games-test
        target: dos-games-test
        executable: true
    actions:
      onDeploy:
        after:
          - cmd: ./dos-games-test -test.v
```
