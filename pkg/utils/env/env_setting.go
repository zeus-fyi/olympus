package env

var Str string

type Env interface {
	Dev() string
	Staging() string
	Production() string
	Local() string
}

func SetEnvParam(e Env) string {
	switch Str {
	case "dev", "development":
		return e.Dev()
	case "staging":
		return e.Staging()
	case "production":
		return e.Production()
	default:
		return e.Local()
	}
}
