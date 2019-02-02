package functions

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func Request(res http.ResponseWriter, req *http.Request, target interface{}) error {
	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		Response(res, req, 400, []byte("{}"))
		return err
	}

	err = json.Unmarshal(body, &target)
	if err != nil {
		Response(res, req, 400, []byte("{}"))
		return err
	}
	return nil
}