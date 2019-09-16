# requests
a go http client for human

## Examples

```go
import "github.com/jacexh/requests"


func main() {
	_, _, err := requests.Get("https://cn.bing.com/search", requests.Params{Query: Any{"q": "golang"}}, nil)
	if err != nil {
		panic(err)
	}

}

```