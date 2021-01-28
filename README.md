# cueflow repro cases

## usage

```sh
go run main.go ./repro/simple
```

## setup

tasks are defined as structures containing an `input` and `output` field.
Upon running, `output` will be set to (pseudo code) `fmt.Sprintf("from %s: %s",
task.Path(), input)`

## repro cases

- :white_check_mark: **nodeps**: `A` and `B` tasks, no dependencies
- :x: **simple**: Simple `B` to `A` direct dependency
- :x: **interpolation**: `B` to `A` dependency through string interpolation
- :x: **unmarshal**: `B` to `A` indirect dependency via a `Call` expression (`json.Unmarshal()`)
