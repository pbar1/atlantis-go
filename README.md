# Atlantis SDK for Go

A Go library for creating [custom `run` commands][1] for [Atlantis][2].

```go
import atlantis "github.com/pbar1/atlantis-go"
```

## Usage

```go
package main

import (
	"log"

	atlantis "github.com/pbar1/atlantis-go"
)

func main() {
	step, err := atlantis.NewRunStep()
	if err != nil {
		log.Fatal(err)
  }
	log.Printf("Pull request number: %d\n", step.PullNum)
	log.Printf("Terraform plan file: %s\n", step.Planfile)
}
```

[1]: https://www.runatlantis.io/docs/custom-workflows.html#reference
[2]: https://www.runatlantis.io

kimchi tofu soup mild, pork
"", beef
