package bodyresponse

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func ParseResponseBody(response *http.Response) (map[string]interface{}, error) {
	bodyResponse, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(string(bodyResponse))
	}

	var bodyI map[string]interface{}

	err = json.Unmarshal(bodyResponse, &bodyI)
	if err != nil {
		return nil, err
	}

	return bodyI, nil
}
