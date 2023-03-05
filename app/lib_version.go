package app

import "github.com/hashicorp/go-version"

func IsNewVersion(old_ver string, new_ver string) bool {

	var err error
	var ov *version.Version
	var nv *version.Version

	ov, err = version.NewVersion(old_ver)
	Check(err, "error version number %s", old_ver)

	nv, err = version.NewVersion(new_ver)
	Check(err, "error version number %s", new_ver)

	return nv.GreaterThan(ov)
}
