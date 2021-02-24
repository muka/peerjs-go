package peer

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

//NewAPI initiate a new API client
func NewAPI(opts Options) API {
	return API{
		opts: opts,
		log:  createLogger("api", opts.Debug),
	}
}

//API wrap calls to API server
type API struct {
	opts Options
	log  *logrus.Entry
}

func (a *API) buildURL(method string) string {
	proto := "http"
	if a.opts.Secure {
		proto = "https"
	}

	path := a.opts.Path
	if path == "/" {
		path = ""
	}

	return fmt.Sprintf(
		"%s://%s:%d%s/%s/%s?ts=%d%d",
		proto,
		a.opts.Host,
		a.opts.Port,
		path,
		a.opts.Key,
		method,
		time.Now().UnixNano(),
		rand.Int(),
	)
}

func (a *API) req(method string) ([]byte, error) {
	uri := a.buildURL(method)
	resp, err := http.Get(uri)
	if err != nil {
		return []byte{}, err
	}
	if resp.StatusCode >= 400 {
		return []byte{}, fmt.Errorf("Request %s failed: %s", uri, resp.Status)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

//RetrieveID retrieve a ID
func (a *API) RetrieveID() ([]byte, error) {
	return a.req("id")
}

//ListAllPeers return the list of available peers
func (a *API) ListAllPeers() ([]byte, error) {
	return a.req("peers")
}
