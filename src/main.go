package main
import (
	"fmt"
	"flag"
	"os"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func main() {
	cur0Ptr := flag.String("cur0", "btc", "first and last currency")
	//amount0Ptr := flag.float("amount0", 100, "starting amount")
	flag.Parse()
	var curs [3]string
	curs[0] = *cur0Ptr
	pairs, minAmount, fee := getPairsArrays()
	var amount [4]float64
	amount[0] = 100 //*amount0Ptr
	index := [3]int{0, 0, 0}
	// Поиск первой промежуточной валюты
	for _, pair := range pairs {
		if ((string(pair[len(pair)-len(curs[0]):len(pair)])) == curs[0]) {
			curs[1] = pair[:len(pair)-4]
			// Поиск второй промежуточной валюты
			index[1] = 0
			for _, pair_2 := range pairs {
				if (len(pair_2) > len(curs[1])) {
					if ((string(pair_2[:len(curs[1])])) == curs[1]) {
						curs[2] = pair_2[len(curs[1])+1:len(pair_2)]
						// Возвращение к первой валюте
						index[2] = 0
						for _, pair_3 := range pairs {
							if (len(pair_3) > len(curs[2])) {
								if pair_3[:len(curs[2])] == curs[2] {
									if pair_3[len(curs[2])+1:len(pair_3)] == curs[0] {
										trans_1, err := getActiveOrders(curs[1] + "_" + curs[0], 1)
										time.Sleep(700 * time.Millisecond)
										trans_2, err := getActiveOrders(curs[1] + "_" + curs[2], 1)
										time.Sleep(700 * time.Millisecond)
										trans_3, err := getActiveOrders(curs[2] + "_" + curs[0], 1)
										time.Sleep(700 * time.Millisecond)
										if err == nil {
											if len(trans_1.Asks) != 0 {
												if len(trans_2.Bids) != 0 {
													if len(trans_3.Bids) != 0 {
														amount[0] = amount[0] - ((amount[0] / 100) * fee[index[0]])
														if amount[0] >= minAmount[index[0]] {
															amount[1] = amount[0] / trans_1.Asks[0][0]
															amount[1] = amount[1] - ((amount[1] / 100) * fee[index[1]])
															if amount[1] >= minAmount[index[1]] {
																amount[2] = amount[1] * trans_2.Bids[0][0]
																amount[2] = amount[2] - ((amount[2] / 100) * fee[index[2]])
																if amount[2] >= minAmount[index[2]] {
																	amount[3] = amount[2] * trans_3.Bids[0][0]
																	if (amount[3] - amount[0]) > 0 {
																		fmt.Println(curs[0])
																		fmt.Println(curs[1])
																		fmt.Println(curs[2])
																		fmt.Println(curs[0])
																		fmt.Println(amount[3] - amount[0])
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
						index[2] = index[2] + 1
						}
					}
				}
			index[1] = index[1] + 1
			}
		}
	index[0] = index[0] + 1
	}
}

type AB struct {
	Asks [][]float64 `json:"asks"`
	Bids [][]float64 `json:"bids"`
}

func getPairsArrays() ([]string, []float64, []float64) {
	url := "https://yobit.net/api/3/info"
	i := 0
	info := make(map[string]interface{})
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
	err = json.Unmarshal(b, &info)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	pairsI := info["pairs"].(map[string]interface{})
	pairs := make([]string, len(pairsI))
	minAmount := make([]float64, len(pairsI))
	fee := make([]float64, len(pairsI))
	for temp, _ := range pairsI {
		pairs[i] = temp
		temp_pair := pairsI[temp].(map[string]interface{})
		minAmount[i] = temp_pair["min_amount"].(float64)
		fee[i] = temp_pair["fee"].(float64)
		i++
	}
	return pairs, minAmount, fee
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
		os.Exit(1)
	}
	json.Unmarshal(b, &ab)
	return &ab, err
}
