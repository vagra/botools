package app

import (
	"errors"
	"fmt"
	"strconv"
)

func Welcome() {
	println(WELCOME)

	var input string
	var step int

	for {
		_, err := fmt.Scanln(&input)
		if err != nil {
			println("请输入")
			continue
		}

		step, err = strconv.Atoi(input)
		if err != nil {
			println("请输入数字")
			continue
		}

		if step < 0 || step > 200 {
			println("请输入前面列举的有效数字")
			continue
		}

		err = Run(step)
		if err != nil {
			println(err.Error())
			continue
		}

		break
	}

	WaitExit(0)
}

func Run(step int) error {

	var err error = nil

	switch step {
	case 1:
		err = InitDBs()
	case 2:
		err = GetTree()
	case 3:
		err = GetSize()
	case 4:
		err = CheckSum()
	case 5:
		err = VirTree()
	case 6:
		err = SyncDB2VDB()
	case 7:
		err = SyncReal2DB()
	// case 101:
	// err = TrimIDs()
	case 102:
		err = ModPaths()
	case 103:
		err = MoveErrors()
	case 200:
		err = MigrateDB()
	case 0:
		Exit(0)
	default:
		err = errors.New("请输入前面列举的有效数字")
	}

	return err
}
