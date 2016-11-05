package teamcity

import(
	"io"
	"errors"
	"net/http"
	"os"
)


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

func (this *Teamcity) DownloadArtifacts(buildConf string, tmp *os.File) (error) {
	// make the http request
	resp, err := this.request("repository/downloadAll/" + buildConf + "/.lastSuccessful/")
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	// copy content of the zip to the tempfile
	_,err = io.Copy(tmp, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
