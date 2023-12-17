package main

func cancleAccount(token string, id string) error {

	if _, err := request([]byte(``), token, id, "cancle"); err != nil {
		panic(err)
	}

	return nil
}
