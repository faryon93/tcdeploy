package teamcity

import (
	"net/http"
	"io/ioutil"
	"encoding/xml"
)


// ----------------------------------------------------------------------------------
//  types
// ----------------------------------------------------------------------------------

type Teamcity struct {
	Url string
	Username string
	Password string

	client *http.Client
}


// ----------------------------------------------------------------------------------
//  constructors
// ----------------------------------------------------------------------------------

func New(url, username, password string) (*Teamcity) {
	return &Teamcity{
		Url: url,
		Username: username,
		Password: password,

		client: &http.Client{},
	}
}


// ----------------------------------------------------------------------------------
//  private members
// ----------------------------------------------------------------------------------

func (this *Teamcity) request(resource string) (*http.Response, error) {
	// make the http request
	req, err := http.NewRequest("GET", this.Url + "/httpAuth/" + resource, nil)
	if err != nil {
		return nil, err
	}

	// setup authentication and performe the request
	req.SetBasicAuth(this.Username, this.Password)
	return this.client.Do(req)
}

func (this *Teamcity) getXml(resource string, v interface{}) (error) {
	resp, err := this.request(resource)
	if err != nil {
		return err
	}

	// read the whole body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return xml.Unmarshal(body, &v)
}