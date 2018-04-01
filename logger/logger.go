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
	if line == "" {
		return
	}
	tm := time.Now()
	toLog := log.SuccessLog + tm.Format("Mon Jan 2 15:04:05 -0700 MST 2006") + " > " + line + "\n"
	log.SuccessLog = toLog
}

//ReportError  Function to write to the errors log
func (log *Logger) ReportError(line string) {
	if line == "" {
		return
	}
	tm := time.Now()
	toLog := log.ErrorsLog + tm.Format("Mon Jan 2 15:04:05 -0700 MST 2006") + " > " + line + "\n"
	log.ErrorsLog = toLog
}

//ReportWarning  Function to write to the warnings log
func (log *Logger) ReportWarning(line string) {
	if line == "" {
		return
	}
	tm := time.Now()
	toLog := log.WarningsLog + tm.Format("Mon Jan 2 15:04:05 -0700 MST 2006") + " > " + line + "\n"
	log.WarningsLog = toLog
}

//WriteLog  Function to write the files in the log directory
func (log *Logger) WriteLog() {

	if _, err := os.Stat("Cooking/log"); err != nil {
		os.Mkdir("Cooking", 0755)
		os.Mkdir("Cooking/log", 0755)
	}

	suffix := "======================================================="
	suffix += suffix + "\n"
	log.SuccessLog += suffix
	log.ErrorsLog += suffix
	log.WarningsLog += suffix

	file, _ := os.OpenFile("Cooking/log/build.success",
		os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	file.Write([]byte(log.SuccessLog))

	file, _ = os.OpenFile("Cooking/log/build.errors",
		os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	file.Write([]byte(log.ErrorsLog))

	file, _ = os.OpenFile("Cooking/log/build.warnings",
		os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	file.Write([]byte(log.WarningsLog))

	file.Close()

}

//NewLogger  Function for creating a new logger object
func NewLogger() Logger {
	tm := time.Now()
	suffix := " BUILD =================================================\n\n"
	log := Logger{
		SuccessLog:  tm.String() + suffix,
		ErrorsLog:   tm.String() + suffix,
		WarningsLog: tm.String() + suffix,
	}

	return log
}
