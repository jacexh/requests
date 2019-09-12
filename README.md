# requests
a go http client for human

## Examples

```go
import "github.com/jacexh/requests"


func main() {
	_, _, err := requests.Get("https://cn.bing.com/search", requests.Parameters{Query: map[string]string{"q": "golang"}}, nil)
	if err != nil {
		panic(err)
	}

}

```