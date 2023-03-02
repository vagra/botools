package app

import "fmt"

type Roots struct {
	disks_root  string
	errors_root string
	dups_root   string
	virs_root   string
}

func (r Roots) Tuple() string {
	return fmt.Sprintln("[roots]") +
		fmt.Sprintf("disks-root = %s\n", r.disks_root) +
		fmt.Sprintf("errors-root = %s\n", r.errors_root) +
		fmt.Sprintf("dups-root = %s\n", r.dups_root) +
		fmt.Sprintf("virs-root = %s", r.virs_root)
}
