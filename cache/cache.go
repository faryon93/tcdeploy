package cache

import (
	"os"
	"encoding/json"
	"io/ioutil"
)


// ----------------------------------------------------------------------------------
//  types
// ----------------------------------------------------------------------------------

type Cache struct {
	LastBuildNumber int `json:"build_number"`
	Files []string `json:"files"`
    Dirs []string `json:"dirs"`
}


// ----------------------------------------------------------------------------------
//  public functions
// ----------------------------------------------------------------------------------

func Load(path string) (*Cache, error) {
	raw, err := ioutil.ReadFile(path)
    if err != nil {
    	if os.IsNotExist(err) {
    		return nil, nil
    	} else {
        	return nil, err
        }
    }

    var c Cache
    json.Unmarshal(raw, &c)
    return &c, nil
}


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

func (this *Cache) Save(path string) (error) {
	bytes, err := json.Marshal(this)
    if err != nil {
        return err
    }

    return ioutil.WriteFile(path, bytes, 0644)
}