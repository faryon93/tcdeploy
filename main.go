package main

import (
	"fmt"
	"log"
	"flag"
	"os"
	"path/filepath"
    "strings"
    "sync"
    "io"
    "io/ioutil"
   	"archive/zip"
   	"time"

	"github.com/faryon93/tcdeploy/teamcity"
	"github.com/faryon93/tcdeploy/config"
	"github.com/faryon93/tcdeploy/cache"
)


// ----------------------------------------------------------------------------------
//  constants
// ----------------------------------------------------------------------------------

const (
	CONFIG_FILE_NAME = "Deployfile"
	DEPLOYCACHE_FILE = ".deploycache"
	ARTIFACT_TEMP_DIR = "/tmp"
	ARTIFACT_TEMP_PREFIX = "tcdeploy"
)


// ----------------------------------------------------------------------------------
//  global variables
// ----------------------------------------------------------------------------------

var processors sync.WaitGroup = sync.WaitGroup{}


// ----------------------------------------------------------------------------------
//  application entry
// ----------------------------------------------------------------------------------

func main() {
    // parse command line args
    flag.Parse()

    // make sure a command is supplied
    if len(flag.Args()) < 1 {
        fmt.Println("usage: tcdeploy watch-dir")
        os.Exit(-1)
    }

    for {
	    // recursively check the watch directory
	    // if any Deployfiles exit
	    filepath.Walk(flag.Args()[0], func(path string, f os.FileInfo, err error) error {
	        if f != nil && !f.IsDir() && strings.HasSuffix(f.Name(), CONFIG_FILE_NAME) {
	        	// load the Deployfile
	            conf, err := config.Load(path)
	            if err != nil {
	            	log.Println("failed to load Deployfile:", err.Error())
	            	return err
	            }

	            // process the Deployfiles in paralell
	            processors.Add(1)
	            go process(*conf)
	        }
	        return nil
	    })

	    // wait until all Deployment files have been processed
	    processors.Wait()

	    // sleep until the next cycle
	    time.Sleep(60 * time.Second)
    }
}


// ----------------------------------------------------------------------------------
//  private functions
// ----------------------------------------------------------------------------------

func process(config config.Config) {
	defer processors.Done()

	// some metadata
	dir := filepath.Dir(config.Path)
	cacheFile := filepath.Join(dir, DEPLOYCACHE_FILE)

	// load the cachefile
	buildCache, err := cache.Load(cacheFile)
	if err != nil {
		log.Println("failed to load cache file:", err.Error())
	}

	// check the last build for the build configuration
	// in order to determen if we have to update the deployment dir
	tc := teamcity.New(config.TcUrl, config.TcUser, config.TcPassword)
	builds, err := tc.GetBuilds(config.TcBuildConfId, 1)
	if err != nil {
		fmt.Println(err)
		return
	}

	// skip this directory, no successfull build is present
	if len(builds) == 0 {
		return
	}

	// the build number has increased or no deployment
	// has been executed yes -> we should do it now
	if buildCache == nil || buildCache.LastBuildNumber < builds[0].Number {
		log.Printf("deploying directory %s (#%d)\n", dir, builds[0].Number)

		// we need to clean the deployment directory
		if buildCache != nil {
			// TODO: 
			log.Println("purging currently deployed files in dir", dir)
		}

		// create a temporary file to download the artifact zip to
		tmp, err := ioutil.TempFile(ARTIFACT_TEMP_DIR, ARTIFACT_TEMP_PREFIX)
		if err != nil {
			log.Println("failed to create temporary artifact file in", dir + ":", err.Error())
			return
		}
		defer os.Remove(tmp.Name())

		// download the artifacts to a tempfile
		log.Println("downloading artifact file to", tmp.Name())
		err = tc.DownloadArtifacts(config.TcBuildConfId, tmp)
		if err != nil {
			log.Println("failed to download artifact file:", err.Error())
			return
		}

		// extract the downloaded zip archive
		log.Println("extracting artifact file to", dir)
		paths, err := unzip(tmp.Name(), dir)
		if err != nil {
			log.Println("failed to extract artifact file:", err.Error())
			return
		}

		// setup the cache entry
		buildCache = &cache.Cache{
			LastBuildNumber: builds[0].Number,
			Paths: paths,
		}

		// save the cachefile with the new infos
		err = buildCache.Save(cacheFile)
		if err != nil {
			log.Println("failed to save cachefile:", err.Error())
			return
		}

		log.Printf("successfull deployed build #%d to %s", builds[0].Number, dir)
	}
}


// ----------------------------------------------------------------------------------
//  helper functions
// ----------------------------------------------------------------------------------

// @see http://stackoverflow.com/questions/20357223/easy-way-to-unzip-file-with-golang
func unzip(src, dest string) ([]string, error) {
    r, err := zip.OpenReader(src)
    if err != nil {
        return nil, err
    }
    defer func() {
        if err := r.Close(); err != nil {
            panic(err)
        }
    }()

    os.MkdirAll(dest, 0755)

    paths := make([]string, 0)

    // Closure to address file descriptors issue with all the deferred .Close() methods
    extractAndWriteFile := func(f *zip.File) error {
        rc, err := f.Open()
        if err != nil {
            return err
        }
        defer func() {
            if err := rc.Close(); err != nil {
                panic(err)
            }
        }()

        path := filepath.Join(dest, f.Name)
        paths = append(paths, path)

        if f.FileInfo().IsDir() {
            os.MkdirAll(path, 0755)

        } else {
            os.MkdirAll(filepath.Dir(path), 0755)
            f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
            if err != nil {
                return err
            }
            defer func() {
                if err := f.Close(); err != nil {
                    panic(err)
                }
            }()

            _, err = io.Copy(f, rc)
            if err != nil {
                return err
            }
        }
        return nil
    }

    for _, f := range r.File {
        err := extractAndWriteFile(f)
        if err != nil {
            return nil, err
        }
    }

    return paths, nil
}