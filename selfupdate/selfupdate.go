// Update protocol:
//
//	GET hk.heroku.com/hk/linux-amd64.json
//
//	200 ok
//	{
//	    "Version": "2",
//	    "Sha256": "..." // base64
//	}
//
// then
//
//	GET hkpatch.s3.amazonaws.com/hk/1/2/linux-amd64
//
//	200 ok
//	[bsdiff data]
//
// or
//
//	GET hkdist.s3.amazonaws.com/hk/2/linux-amd64.gz
//
//	200 ok
//	[gzipped executable data]
package selfupdate

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
)

const plat = runtime.GOOS + "-" + runtime.GOARCH

var ErrHashMismatch = errors.New("new file hash mismatch after patch")
var up = NewUpdate()
var defaultHTTPRequester = HTTPRequester{}

// Updater is the configuration and runtime data for doing an update.
//
// Note that ApiURL, BinURL and DiffURL should have the same value if all files are available at the same location.
//
// Example:
//
//	updater := &selfupdate.Updater{
//		CurrentVersion: version,
//		ApiURL:         "http://updates.yourdomain.com/",
//		BinURL:         "http://updates.yourdownmain.com/",
//		DiffURL:        "http://updates.yourdomain.com/",
//		Dir:            "update/",
//		CmdName:        "myapp", // app name
//	}
//	if updater != nil {
//		go updater.BackgroundRun()
//	}
type Updater struct {
	CurrentVersion string    // Currently running version.
	ApiURL         string    // Base URL for API requests (json files).
	BinURL         string    // Base URL for full binary downloads.
	Requester      Requester //Optional parameter to override existing http request handler
	Info           struct {
		Version string
		Sha256  []byte
	}
}

func (u *Updater) GetExecRelativeDir(dir string) string {
	filename, _ := os.Executable()
	path := filepath.Join(filepath.Dir(filename), dir)
	return path
}

// BackgroundRun starts the update check and apply cycle.
func (u *Updater) BackgroundRun() error {
	if u.WantUpdate() {
		if err := up.CanUpdate(); err != nil {
			// fail
			return err
		}

		// TODO(bgentry): logger isn't on Windows. Replace w/ proper error reports.
		if err := u.Update(); err != nil {
			return err
		}
	}
	return nil
}

// WantUpdate returns boolean designating if an update is desired
func (u *Updater) WantUpdate() bool {
	return u.CurrentVersion != "dev"
}

// UpdateAvailable checks if update is available and returns version
func (u *Updater) UpdateAvailable() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", err
	}
	old, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer old.Close()

	err = u.fetchInfo()
	if err != nil {
		return "", err
	}
	if u.Info.Version == u.CurrentVersion {
		return "", nil
	} else {
		return u.Info.Version, nil
	}
}

// Update initiates the self update process
func (u *Updater) Update() error {
	path, err := os.Executable()
	if err != nil {
		return err
	}
	old, err := os.Open(path)
	if err != nil {
		return err
	}
	defer old.Close()

	err = u.fetchInfo()
	if err != nil {
		return err
	}
	if u.Info.Version == u.CurrentVersion {
		return nil
	}

	bin, err := u.fetchAndVerifyFullBin()
	if err != nil {
		if err == ErrHashMismatch {
			log.Println("update: hash mismatch from full binary")
		} else {
			log.Println("update: fetching full binary,", err)
		}
		return err
	}

	// close the old binary before installing because on windows
	// it can't be renamed if a handle to the file is still open
	old.Close()

	err, errRecover := up.FromStream(bytes.NewBuffer(bin))
	if errRecover != nil {
		return fmt.Errorf("update and recovery errors: %q %q", err, errRecover)
	}
	if err != nil {
		return err
	}
	return nil
}

func (u *Updater) fetchInfo() error {
	r, err := u.fetch(u.ApiURL + "/" + url.QueryEscape(plat) + ".json")
	if err != nil {
		return err
	}
	defer r.Close()
	err = json.NewDecoder(r).Decode(&u.Info)
	if err != nil {
		return err
	}
	if len(u.Info.Sha256) != sha256.Size {
		return errors.New("bad cmd hash in info")
	}
	return nil
}

func (u *Updater) fetchAndVerifyFullBin() ([]byte, error) {
	bin, err := u.fetchBin()
	if err != nil {
		return nil, err
	}
	verified := verifySha(bin, u.Info.Sha256)
	if !verified {
		return nil, ErrHashMismatch
	}
	return bin, nil
}

func (u *Updater) fetchBin() ([]byte, error) {
	r, err := u.fetch(u.BinURL + "/" + url.QueryEscape(plat) + ".gz")
	if err != nil {
		return nil, err
	}
	defer r.Close()
	buf := new(bytes.Buffer)
	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(buf, gz); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (u *Updater) fetch(url string) (io.ReadCloser, error) {
	if u.Requester == nil {
		return defaultHTTPRequester.Fetch(url)
	}

	readCloser, err := u.Requester.Fetch(url)
	if err != nil {
		return nil, err
	}

	if readCloser == nil {
		return nil, fmt.Errorf("Fetch was expected to return non-nil ReadCloser")
	}

	return readCloser, nil
}

func verifySha(bin []byte, sha []byte) bool {
	h := sha256.New()
	h.Write(bin)
	return bytes.Equal(h.Sum(nil), sha)
}
