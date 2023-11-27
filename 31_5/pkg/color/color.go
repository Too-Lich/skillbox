package color

func Red() string {
	return "\033[31m"
}

func Green() string {
	return "\033[32m"
}

func Yellow() string {
	return "\033[33m"
}

func Rst() string { //сброс цвета
	return "\033[0m"
}
