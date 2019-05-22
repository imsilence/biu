module biu

go 1.12

replace golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190513172903-22d7a77e9e5f

replace golang.org/x/sys => github.com/golang/sys v0.0.0-20190522044717-8097e1b27ff5

replace golang.org/x/net => github.com/golang/net v0.0.0-20190520210107-018c4d40a106

replace golang.org/x/text => github.com/golang/text v0.3.2

replace golang.org/x/tools => github.com/golang/tools v0.0.0-20190521203540-521d6ed310dd

replace golang.org/x/sync => github.com/golang/sync v0.0.0-20190423024810-112230192c58

require (
	github.com/apparentlymart/go-cidr v1.0.0
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0
	github.com/levigross/grequests v0.0.0-20190130132859-37c80f76a0da
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.4
)
