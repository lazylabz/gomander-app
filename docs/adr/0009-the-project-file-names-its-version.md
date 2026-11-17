# The project file names its version, and a file without one is refused

The file a Project is exported to carries a top-level `version`, and
`ProjectFile.Read` picks the codec for that number
(`internal/project/infrastructure/projectfile.go`). A number with no codec,
a missing `version` included, is refused with `unsupported project file version`.

Every released Gomander wrote the field, so a file without it was not written by
Gomander:

| Releases        | Writes       | Encoder                                                  |
| --------------- | ------------ | -------------------------------------------------------- |
| v1.0.0          | `version: 1` | `internal/app/importer.go`                               |
| v1.1.0 – v1.8.0 | `version: 1` | `internal/project/application/usecases/exportproject.go` |
| v1.9.0 – v2.1.0 | `version: 1` | `internal/project/infrastructure/projectfilev1.go`       |
| v2.1.1 onwards  | `version: 2` | `internal/project/infrastructure/projectfilev2.go`       |

Version 2 adds each command's link and error patterns, which version 1 never
carried.

## Rejected alternatives

- **Reading a missing `version` as 1.** It would accept files no Gomander produced
  and guess at their shape, while no file a user exported needs the leniency.
- **Changing a codec in place when the format grows.** A file exported by an older
  Gomander must keep importing exactly as it did, so a new shape is a new codec and
  a new number; `projectfilev1.go` stays byte for byte.

## Consequences

A format change adds `projectfilevN.go`, a `case N` in `Read`, and points `Write`
at it, and the table above gains a row. Older codecs are decode-only.
