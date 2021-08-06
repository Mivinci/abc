# Webkit

An implementation of HTTP router for learning purposes.


## Usage

```bash
go get -u github.com/mivinci/webkit/v2
```

## Example

```go
package main

import (
    "net/http"
    "strconv"

    "github.com/mivinci/webkit/v2"
)

func main() {
    h := webkit.New()

    h.Get("/:id", func(c webkit.Ctx) error {
        id, err := c.Params().Int("id")
        if err != nil {
            return webkit.HTTPError(http.StatusBadRequest)
    }
        return c.Text(http.StatusOK, strconv.Itoa(id))
    })

    h.Run(":8080")
}
```
