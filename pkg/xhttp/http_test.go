package xhttp

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMyClient_DoGet(t *testing.T) {
	c := NewClient(time.Minute)

	header := map[string]string{
		"X-Forward-IP": "127.0.0.1",
		"X-Trace-ID":   "1234",
	}

	Query := map[string]string{
		"a": fmt.Sprint(1),
		"b": "c",
	}

	resp, err := c.DoGet(context.TODO(), "", header, Query)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	defer resp.Body.Close()

	fmt.Println(string(body))
}

func TestMyClient_DoPOST(t *testing.T) {
	c := NewClient(time.Minute)

	header := map[string]string{
		"X-Forward-IP": "127.0.0.1",
		"X-Trace-ID":   "1234",
	}

	req := struct {
		Name string `json:"name"`
		UID  int64  `json:"uid"`
	}{
		UID:  10,
		Name: "foo",
	}

	resp, err := c.DoPost(context.TODO(), "", header, req)

	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	defer resp.Body.Close()

	fmt.Println(string(body))
}
