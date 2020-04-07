package main
import (
	"fmt"
	"os"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func main() {

	url := "https://yobit.net/api/3/info"

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	b, err := ioutil.ReadAll(resp.Body)
	
	resp.Body.Close()
	
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	info := make(map[string]interface{})

	err = json.Unmarshal(b, &info)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pairs := info["pairs"].(map[string]interface{})

	k := make([]string, len(pairs))

	i := 0

	for s, _ := range pairs {
		k[i] = s
		i++
	}

	fmt.Printf("%#v\n", k)

}


