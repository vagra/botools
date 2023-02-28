package app

func InitMaps() {
	g_map_dirs = make(map[string]map[string]*Dir)
	g_map_files = make(map[string]map[string]*File)

	g_dirs_counter = make(map[string]*int64)
	g_files_counter = make(map[string]*int64)
}

func InitMap(disk_name string) {
	g_map_dirs[disk_name] = make(map[string]*Dir)
	g_map_files[disk_name] = make(map[string]*File)

	var dirs_counter int64 = 0
	var files_counter int64 = 0

	g_dirs_counter[disk_name] = &dirs_counter
	g_files_counter[disk_name] = &files_counter
}

func CompDirsPath(disk_name string) {
	var count int64 = 0
	for {
		count = GetNoPathDirsCount(disk_name)
		if count == 0 {
			break
		}

		DirsAddParentPath(disk_name)
	}
}

func CompFilesPath(disk_name string) {
	FilesAddParentPath(disk_name)
}

func DirsAddParentPath(disk_name string) {

	for _, dir := range g_map_dirs[disk_name] {
		if dir.parent_id == "0" {
			continue
		}

		if len(dir.path) > 0 {
			continue
		}

		parent_path := g_map_dirs[disk_name][dir.parent_id].path

		if len(parent_path) == 0 {
			continue
		}

		dir.path = parent_path + "/" + dir.name
	}
}

func FilesAddParentPath(disk_name string) {

	for _, file := range g_map_files[disk_name] {

		if len(file.path) > 0 {
			continue
		}

		parent_path := g_map_dirs[disk_name][file.parent_id].path

		if len(parent_path) == 0 {
			continue
		}

		file.path = parent_path + "/" + file.name
	}
}

func GetNoPathDirsCount(disk_name string) int64 {

	var count int64 = 0
	for _, dir := range g_map_dirs[disk_name] {
		if len(dir.path) == 0 {
			count++
		}
	}

	return count
}
