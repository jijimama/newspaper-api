package presenter

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func PrettyPrintStructToJson(data interface{}) {
	var buffer bytes.Buffer
	enc := json.NewEncoder(&buffer)
	enc.SetIndent("", "    ")
	if err := enc.Encode(data); err != nil {
		fmt.Println(err)
	}
	fmt.Println(buffer.String())
}
