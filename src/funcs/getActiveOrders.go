package main
import (
	"fmt"
	"net/http"
	"os"
	"encoding/json"
	"io/ioutil"	
	"strconv"
)

type AB struct {
	Asks [][]float64 `json:"asks"`
	Bids [][]float64 `json:"bids"`
}

func main() {

	pair := "ltc_btc"

	limit := 3 

	url := "https://yobit.net/api/3/depth/" + pair + "?limit=" + strconv.FormatInt(int64(limit), 10) 

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(resp.Body)	

	resp.Body.Close()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var result map[string]interface{}

	err = json.Unmarshal(body, &result)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	result = result[pair].(map[string]interface{})

	b, err := json.Marshal(result)

	var ab AB

	json.Unmarshal(b, &ab)

	fmt.Println(ab.Asks)

}
