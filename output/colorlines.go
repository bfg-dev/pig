package output

import (
	"log"
)

// Fatal - log fatal
func Fatal(v ...interface{}) {
	log.Fatal(red(v))
}

// Fatalf - log fatal
func Fatalf(format string, v ...interface{}) {
	log.Fatal(redf(format, v))
}

// Info1 - log info
func Info1(v ...interface{}) {
	log.Println(yellow(v))
}

// Infof1 - log info
func Infof1(format string, v ...interface{}) {
	log.Println(yellowf(format, v))
}

// Info2 - log info
func Info2(v ...interface{}) {
	log.Println(cyan(v))
}

// Infof2 - log info
func Infof2(format string, v ...interface{}) {
	log.Println(cyanf(format, v))
}

// OK - log info
func OK(v ...interface{}) {
	log.Println(green(v))
}

// OKf - log info
func OKf(format string, v ...interface{}) {
	log.Println(greenf(format, v))
}
