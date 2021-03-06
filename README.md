# tcdeploy - teamcity deployment service

[![Travis CI](https://travis-ci.org/faryon93/tcdeploy.svg?branch=master)](https://travis-ci.org/faryon93/tcdeploy)

---

Tcdeploy offers a simple method to deploy artifacts from TeamCity build servers to your target server. The build server is queried cyclically for a new successfull build in order to download and extract all artifacts to the target directory. We support deploying multiple build configurations to multiple directories.


## Usage

The only command-line option, which is obligatory, is the directory which should be search recursivly by tcdeploy:
```sh
$ /usr/sbin/tcdeploy /var/www
```

Tcdeploy is designed as a daemon, so a systemd service file is supplied within this repository.

## Sample Configuration
Place a **Deployfile** file in a directory within the search path of tcdeploy. The artifacts will be deployed in the directory containing the **Deployfile**.

```ini
provider = "tc"
tc_url = "http://teamcity.local"
tc_build_conf = "MyBuildConfig"
tc_user = "deploy"
tc_password = "deploy"

```

## Security Considerations

If the directory, where the Deployfile is located, will be published to the general public (e.g. via http server), then you should restrict the access to the following files:

- Deployfile
- .deploycache
