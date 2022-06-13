package helper

import (
	"fmt"
	"io"
)

func CloseConnection(connect io.ReadCloser) {
	err := connect.Close()
	if err != nil {
		fmt.Println("rapid:error close connection:" + err.Error())
	}
}
