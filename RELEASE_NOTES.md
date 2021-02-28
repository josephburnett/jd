# Release Notes

## 1.3

### Features

* Include web tool with commandline binary
* Add Dockerfile

### Fixes

* Statically link release binary
* Convert arrays tests to tdt (wafuwafu13)
* Infer set keys from diff

## 1.2

### Features

* Web tool
* Diff and patch YAML
* -setkeys flag to index objects in a set
* Diff format includes metadata (set, multiset, setkeys) to be self-contained

### Fixes

* Fixes application of patch created when diff is created in set mode (aftab-a)
* Changes jsonNull into a type that unmarshals to null (kevburnsjr)
* Remove panics from patch (kevburnsjr)

## 1.1

### Features

* -set and -mset flag to treat arrays as sets or multisets

## 1.0

* Initial release