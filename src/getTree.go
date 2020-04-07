package main
import (
	"fmt"
	"os"
	"encoding/json"
	"io/ioutil"
	"net/http"
	//"strings"
	"strconv"
	"time"
)

func main() {

	var curs [3]string

	curs[0] = "btc"

	pairs := getPairsArray()

	var money [4]float64

	money[0] = 100

	fee := 0.02

	//i := 0

	// Поиск первой промежуточной валюты
	for _, pair := range pairs {
		if ((string(pair[len(pair)-3:len(pair)])) == curs[0]) {
			curs[1] = pair[:len(pair)-4]
			// Поиск второй промежуточной валюты
			for _, pair_2 := range pairs {
				if (len(pair_2) > len(curs[1])) {
					if ((string(pair_2[:len(curs[1])])) == curs[1]) {
						curs[2] = pair_2[len(curs[1])+1:len(pair_2)]
						// Возвращение к первой валюте
						for _, pair_3 := range pairs {
							if (len(pair_3) > len(curs[2])) {
								if pair_3[:len(curs[2])] == curs[2] {
									if pair_3[len(curs[2])+1:len(pair_3)] == curs[0] {
										fmt.Println("\n")
										trans_1, err := getActiveOrders(curs[1] + "_" + curs[0], 1)
										time.Sleep(700 * time.Millisecond)
										trans_2, err := getActiveOrders(curs[1] + "_" + curs[2], 1)
										time.Sleep(700 * time.Millisecond)
										trans_3, err := getActiveOrders(curs[2] + "_" + curs[0], 1)
										time.Sleep(700 * time.Millisecond)
										if err == nil {
											if len(trans_1.Asks) != 0 {
												money[0] = money[0] - ((money[0] / 100) * fee)
												money[1] = money[0] / trans_1.Asks[0][0]
												if len(trans_2.Bids) != 0 {
													money[1] = money[1] - ((money[1] / 100) * fee)
													money[2] = money[1] * trans_2.Bids[0][0]
													if len(trans_3.Bids) != 0 {
														money[2] = money[2] - ((money[2] / 100) * fee)
														money[3] = money[2] * trans_3.Bids[0][0]
														if (money[3] - money[0]) > 0 {
															fmt.Println(curs[0])
															fmt.Println(curs[1])
															fmt.Println(curs[2])
															fmt.Println(curs[0])
															fmt.Println(money[3] - money[0])
															fmt.Println("\n")
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

type AB struct {
	Asks [][]float64 `json:"asks"`
	Bids [][]float64 `json:"bids"`
}

func getPairsArray() []string {

	url := "https://yobit.net/api/3/info"

	i := 0

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

	for s, _ := range pairs {
		k[i] = s
		i++
	}

	return k
}

func getActiveOrders(pair string, limit int) (*AB, error) {

	url := "https://yobit.net/api/3/depth/" + pair + "?limit=" + strconv.FormatInt(int64(limit), 10)

	var ab AB

	var result map[string]interface{}

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

	err = json.Unmarshal(body, &result)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	result = result[pair].(map[string]interface{})

	b, err := json.Marshal(result)

	if err != nil {
		fmt.Println(err)
	}

	json.Unmarshal(b, &ab)

	return &ab, err
}
