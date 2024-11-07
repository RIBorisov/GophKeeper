package client

import "strconv"

var (
	// ChooseKind is a string using to ask user about kind of saving data.
	ChooseKind = `
Choose kind of data:
[1] - Card
[2] - Text
[3] - Login and password
[4] - Binary
[0] - Return to previous menu`

	// InputAction is a string using to ask user about action.
	InputAction = `Please, input number of action you need [1..6]
----------------------------
[1] - Register
[2] - Log in
[3] - Get data
[4] - Save data
[5] - Exit the application
[6] - Build info
----------------------------`
)

type Action uint

const (
	Register Action = iota + 1
	LogIn
	GetData
	SaveData
	Return
	BuildInfo
)

func (a Action) String() string {
	return strconv.Itoa(int(a))
}
