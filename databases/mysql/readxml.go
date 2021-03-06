package mysql

import (
  "encoding/xml"
  "fmt"
  "os"
)

func ReadXML(file string) *Database {
  xmlFile, err := os.Open(file)
  if nil != err {
    panic(fmt.Sprintf("Error opening file:", err))
  }
  defer xmlFile.Close()

  decoder := xml.NewDecoder(xmlFile)

  var dump Database
  decoder.Decode(&dump)

  return &dump
}
