package client

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
	InputAction = `Please, input number of action you need [1, 2, 3, 4, 9, 0]
----------------------------
[1] - Register
[2] - Log in
[3] - Get data
[4] - Save data
[9] - Build info
[0] - Exit the application
----------------------------`
)

type Action string

const (
	Register  Action = "1"
	LogIn     Action = "2"
	GetData   Action = "3"
	SaveData  Action = "4"
	BuildInfo Action = "9"
	Return    Action = "0"
)

func (a Action) String() string {
	return string(a)
}
