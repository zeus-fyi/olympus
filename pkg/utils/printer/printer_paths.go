package printer

type PrintPath string

func (p PrintPath) Local() string {
	return "artifacts/local/"
}

func (p PrintPath) Dev() string {
	return "artifacts/dev/"
}

func (p PrintPath) Staging() string {
	return "artifacts/staging/"
}

func (p PrintPath) Production() string {
	return "artifacts/production/"
}
