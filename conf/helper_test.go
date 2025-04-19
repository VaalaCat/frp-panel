package conf

import (
	"fmt"
	"net/url"
	"testing"
)

func TestGetRPCConnInfo(t *testing.T) {
	parsedUrl, err := url.Parse("grpc://123123:88/123123")
	if err != nil {
	}

	connInfo := ConnInfo{
		Host:   parsedUrl.Host,
		Scheme: Scheme(parsedUrl.Scheme),
	}

	fmt.Printf("%+v", connInfo)
}
