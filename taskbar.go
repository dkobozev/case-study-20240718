// Try to determine the color of the user's taskbar in Windows 10+.
//
// It doesn't appear that the Windows API provides a dedicated method for
// querying the taskbar color specifically. The following method uses the
// background color as the proxy for taskbar color, which generally matches the
// Light/Dark theme selected by the user, but doesn't work reliably in all
// cases, such as when the user chooses to show the accent color on the taskbar,
// or selects a Custom colorscheme with the Light option for the default Windows
// mode.

package main

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/zzl/go-win32api/v2/win32"
	"github.com/zzl/go-winrtapi/winrt"
)

// define our own struct to represent RGBA color because winrt doesn't provide one
type Color struct {
	A byte
	R byte
	G byte
	B byte
}

func main() {
	winrt.Initialize()

	// activate UISettings
	// https://learn.microsoft.com/en-us/uwp/api/windows.ui.viewmanagement.uisettings
	var p *win32.IInspectable
	hs := winrt.NewHStr("Windows.UI.ViewManagement.UISettings")
	hr := win32.RoActivateInstance(hs.Ptr, &p)
	if win32.FAILED(hr) {
		panic("Failed to activate")
	}

	var uiSettings3 *winrt.IUISettings3
	hr = p.QueryInterface(&winrt.IID_IUISettings3, unsafe.Pointer(&uiSettings3))

	if win32.FAILED(hr) {
		panic("Failed to query")
	}

	background := Color{}
	GetColorValue(uiSettings3, winrt.UIColorType_Background, &background)

	fmt.Printf("Color: #%02x%02x%02x\n", background.R, background.G, background.B)
}

// winrt has a wrapper for GetColorValue, but it doesn't appear to be working correctly
func GetColorValue(uiSettings3 *winrt.IUISettings3, desiredColor winrt.UIColorType, color *Color) {
	syscall.SyscallN(uiSettings3.Vtbl().GetColorValue, uintptr(unsafe.Pointer(uiSettings3)), uintptr(desiredColor), uintptr(unsafe.Pointer(color)))
}
