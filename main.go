package main

import (
  "fmt"
  "github.com/scalia/mysynql/options"
  "os"
)

func main() {
  if !options.Parse() {
    os.Exit(2)
  }

  // Execute as needed.
  opts := &options.ProgramOptions
  if opts.Version {
    fmt.Printf("MySynQL version %s.\n", options.Version)
  } else if "" != opts.DumpFile {
    dump()
  } else if "" != opts.StructureFile {
    restore()
  }
}
