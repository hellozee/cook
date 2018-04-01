package logger

import (
	"os"
	"time"
)

//Logger  Data Structure for holding data which is going to be written
type Logger struct {
	SuccessLog  string
	ErrorsLog   string
	WarningsLog string
}

//ReportSuccess  Function to write to the success log
func (log *Logger) ReportSuccess(line string) {
	tm := time.Now()
	toLog := log.SuccessLog + tm.String() + " > " + line + "\n"
	log.SuccessLog = toLog
}

//ReportError  Function to write to the errors log
func (log *Logger) ReportError(line string) {
	tm := time.Now()
	toLog := log.ErrorsLog + tm.String() + " > " + line + "\n"
	log.ErrorsLog = toLog
}

//ReportWarning  Function to write to the warnings log
func (log *Logger) ReportWarning(line string) {
	tm := time.Now()
	toLog := log.WarningsLog + tm.String() + " > " + line + "\n"
	log.WarningsLog = toLog
}

//WriteLog  Function to write the files in the log directory
func (log *Logger) WriteLog() {

	if _, err := os.Stat("Cooking/log"); err != nil {
		os.Mkdir("Cooking/log", 0755)
	}

	log.SuccessLog += "\n==================================\n"
	log.ErrorsLog += "\n==================================\n"
	log.WarningsLog += "\n==================================\n"

	file, _ := os.OpenFile("Cooking/build.sucess",
		os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	file.Write([]byte(log.SuccessLog))

	file, _ = os.OpenFile("Cooking/build.error",
		os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	file.Write([]byte(log.ErrorsLog))

	file, _ = os.OpenFile("Cooking/build.warnings",
		os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	file.Write([]byte(log.WarningsLog))

	file.Close()

}

//NewLogger  Function for creating a new logger object
func NewLogger() Logger {
	tm := time.Now()
	log := Logger{
		SuccessLog:  tm.String() + " Build =================\n\n\n",
		ErrorsLog:   tm.String() + " Build =================\n\n\n",
		WarningsLog: tm.String() + " Build =================\n\n\n",
	}

	return log
}
