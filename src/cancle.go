package main

func cancleAccount(token string, id string) error {
	_, err := request([]byte(``), token, id, "cancle", "DELETE")
	if err != nil {
		panic(err)
	}
	return nil
}
