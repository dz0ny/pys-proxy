package pys_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"proxy/pys"
	"strings"
	"testing"

	"github.com/Bplotka/go-httpt"
	"github.com/stretchr/testify/require"
)

func TestPostParsing(t *testing.T) {

	data := strings.NewReader("action=pys_api_event&pixel=facebook&event=PageView&data%5Bpage_title%5D=Home&data%5Bpost_type%5D=page&data%5Bpost_id%5D=19&data%5Buser_role%5D=guest&data%5Bplugin%5D=PixelYourSite&data%5Bevent_url%5D=localhost%2F&ids%5B%5D=1234&eventID=zjFSkXLeKe0JkUY3l87kTDD2CibuFGb5vllt&woo_order=&edd_order=")
	request, _ := http.NewRequest(http.MethodPost, "/pys-pixel-proxy", data)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	request.Header.Set("Host", "domain.com")
	request.Header.Set("X-Real-IP", "20.10.30.10")
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.96 Safari/537.36")
	request.Header.Set("X-Requested-With", "XMLHttpRequest")
	response := httptest.NewRecorder()

	pys.Handler(response, request)

	got := response.Body.String()
	require.Contains(t, got, "PageView")

}

func TestFBPost(t *testing.T) {

	s := httpt.NewServer(t)
	s.On(httpt.POST, httpt.AnyPath).Push(func(req *http.Request) (*http.Response, error) {
		body, _ := ioutil.ReadAll(req.Body)
		require.Equal(t, "access_token=token&data=%5B%7B%22test%22%3A%22data%22%7D%5D", string(body))
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
		}, nil
	})
	testClient := s.HTTPClient()
	pys.PostToFB("id", "token", []byte(`{"test":"data"}`), testClient)

}
