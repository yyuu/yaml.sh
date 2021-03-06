/* a utility to get value(s) from YAML */
package main

import (
  "errors"
  "fmt"
  "io/ioutil"
  "os"
  "strconv"
  "strings"
  "gopkg.in/yaml.v1"
)

func main() {
  if len(os.Args) < 2 {
    fmt.Fprintf(os.Stderr, "%s: missing operand\n", os.Args[0])
    os.Exit(1)
  }

  // Read a YAML from stdin
  raw, err := ioutil.ReadAll(os.Stdin)
  if err != nil {
    panic(err)
  }

  data, err := loadYaml(raw)
  if err != nil {
    panic(err)
  }

  switch os.Args[1] {
    case "get":
      doGet(data, os.Args[2:])
    default:
      fmt.Fprintf(os.Stderr, "%s: unknown operation: %s\n", os.Args[1])
      os.Exit(1)
  }
}

func loadYaml(raw []byte) (interface{}, error) {
  var data interface{}
  err := yaml.Unmarshal(raw, &data)
  return data, err
}

func doGet(data interface{}, path []string) {
  if len(path) < 1 {
    fmt.Println(prettyPrint(data))
  } else {
    if dataMap, ok := data.(map[interface{}] interface{}); ok {
      if datum, ok := dataMap[path[0]]; ok {
        doGet(datum, path[1:])
      } else {
        panic(errors.New("not a member: " + path[0]))
      }
    } else {
      panic(errors.New("not a map"))
    }
  }
}

func prettyPrint(data interface{}) string {
  switch x := data.(type) {
    case []interface{}:
      y := make([]string, len(x))
      for i := range x {
        y[i] = strconv.Quote(prettyPrint(x[i]))
      }

      // bash array style
      return fmt.Sprintf("(%s)", strings.Join(y, " "))
    case map[interface{}] interface{}:
      y := make([]string, len(x))
      i := 0
      for key, val := range x {
        y[i] = fmt.Sprintf("%s=%s", key, strconv.Quote(prettyPrint(val)))
        i++
      }

      // shell variable style
      return strings.Join(y, "\n")
    default:
      return fmt.Sprintf("%s", x)
  }
}
