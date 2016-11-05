package teamcity

import (
	"strconv"
)


// ----------------------------------------------------------------------------------
//  types
// ----------------------------------------------------------------------------------

type Build struct {
	Id int `xml:"id,attr"`
	Number int `xml:"number,attr"`
	Status string `xml:"status,attr"`
	State string `xml:"state,attr"`
}

type buildWrapper struct {
	Builds []Build `xml:"build"`
}


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

func (this *Teamcity) GetBuilds(buildConfiguration string, limit int) ([]Build, error) {
	var response buildWrapper
	err := this.getXml("api/buildTypes/" + buildConfiguration + "/builds?status=SUCCESS&count=" + strconv.Itoa(limit), &response)
	if err != nil {
		return nil, err
	}

	return response.Builds, nil
}