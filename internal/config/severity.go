package config

import (
	"fmt"
)

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo	Severity = "info"
)

func (s Severity) IsValid() bool{
	switch s {
		case SeverityError,SeverityWarning,SeverityInfo:
			return true
		default:
			return false
	}
}

func ParseSeverity(val any,defaultSev Severity) (*Severity,error){
	if b,ok := val.(bool);ok{
		if!b{
			return  nil,nil
		}

		return nil,fmt.Errorf("invalid value 'true' - use 'error','warning' or 'info'")
	}

	if s,ok := val.(string);ok{
		sev :=Severity(s)
		if !sev.IsValid(){
			return nil,fmt.Errorf("invalid severity '%s' - valid values: error,warning,info or false to disable",s)
		}
		return &sev,nil
	}

	if val==nil{
		return &defaultSev,nil
	}

	return nil,fmt.Errorf("invalid type - must be sting ('error','warning','info') or false to disable")
}

