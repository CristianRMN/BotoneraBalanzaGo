package model

import (
	"context"
	"os/exec"

	"golang.org/x/sys/windows"
)

//Constantes y variables del programa

var (
	NavegadorCmd       *exec.Cmd
	User32                    = windows.NewLazySystemDLL("user32.dll")
	ProcFindWindowW           = User32.NewProc("FindWindowW")
	ProcSetWindowPos          = User32.NewProc("SetWindowPos")
	AuxIP              string = "192.180.2.44"
	UrlBajarTurnoArray string
	UrlSubirTurnoArray string
	GoroutineCtx       context.Context
	CancelGourutine    context.CancelFunc
)

const (
	SWP_NOZORDER   = 0x0004
	SWP_SHOWWINDOW = 0x0040
)
