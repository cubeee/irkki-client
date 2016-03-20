package log

import (
	"log"
)

var enabled = true

func DisableLogging() {
	enabled = false
}

func EnableLogging() {
	enabled = true
}

func Fatal(v ...interface{}) {
	if enabled {
		log.Fatal(v...)
	}
}

func Fatalf(format string, v ...interface{}) {
	if enabled {
		log.Fatalf(format, v...)
	}
}

func Fatalln(v ...interface{}) {
	if enabled {
		log.Fatalln(v...)
	}
}

func Panic(v ...interface{}) {
	if enabled {
		log.Panic(v...)
	}
}

func Panicf(format string, v ...interface{}) {
	if enabled {
		log.Panicf(format, v...)
	}
}

func Panicln(v ...interface{}) {
	if enabled {
		log.Panicln(v...)
	}
}

func Print(v ...interface{}) {
	if enabled {
		log.Print(v...)
	}
}

func Printf(format string, v ...interface{}) {
	if enabled {
		log.Printf(format, v...)
	}
}

func Println(v ...interface{}) {
	if enabled {
		log.Println(v...)
	}
}
