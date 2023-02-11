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

		if step < 0 || step > STEP_COUNT {
			println("请输入前面列举的有效数字")
			continue
		}

		err = Run(step)
		if err != nil {
			println(err)
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
		err = InitDB()
	case 2:
		err = GetTree()
	case 3:
		err = GetSize()
	case 4:
		err = CheckSum()
	case 5:
		err = VirTree()
	case 0:
		Exit(0)
	default:
		err = errors.New("只能输入 0 到 " + strconv.Itoa(STEP_COUNT) + " 的整数")
	}

	return err
}
