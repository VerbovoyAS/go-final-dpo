package st

import "fmt"

type SMSData struct {
	Country      string
	Bandwidth    string
	ResponseTime string
	Provider     string
}

// ToString преобразует структуру в строку
func (s *SMSData) ToString() string {
	return fmt.Sprintf("%s;%s;%s;%s\n", s.Country, s.Bandwidth, s.ResponseTime, s.Provider)
}
