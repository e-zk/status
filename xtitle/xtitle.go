package xtitle

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

/*func List(X *xgbutil.XUtil) (windows []xproto.Window) {

	// list of clients (windows?)
	clientids, err := ewmh.ClientListGet(X)
	if err != nil {
		//fmt.Printf("Error| %s", err)
		panic(err)
	}

	// for each client...
	for _, clientid := range clientids {
		//fmt.Printf("%#08x\n", clientid)
		windows = append(windows, clientid)
	}

	return windows
}*/

func GetFocus(X *xgb.Conn) (win xproto.Window, err error) {

	setup := xproto.Setup(X)
	rootId := setup.DefaultScreen(X).Root

	aname := "_NET_ACTIVE_WINDOW"
	activeAtom, err := xproto.InternAtom(X, true, uint16(len(aname)),
		aname).Reply()
	if err != nil {
		return
	}

	reply, err := xproto.GetProperty(X, false, rootId, activeAtom.Atom,
		xproto.GetPropertyTypeAny, 0, (1<<32)-1).Reply()
	if err != nil {
		return
	}

	currentId := xproto.Window(xgb.Get32(reply.Value))
	return currentId, nil
}

// Print the window name for the given window...
func GetName(X *xgb.Conn, win xproto.Window) (windowName string, err error) {
	response := new(xproto.GetPropertyReply)

	// get name
	netWmName := "_NET_WM_NAME"
	netWmAtom, err := xproto.InternAtom(X, true, uint16(len(netWmName)),
		netWmName).Reply()
	if err != nil {
		return "", err
	}

	response, err = xproto.GetProperty(X, false, win, netWmAtom.Atom,
		xproto.GetPropertyTypeAny, 0, (1<<32)-1).Reply()
	if err != nil {
		return "", err
	}

	// if there is no _NET_WM_NAME, attempt to get WM_NAME instead...
	if response.Length == 0 {
		wmName := "WM_NAME"
		wmNameAtom, err := xproto.InternAtom(X, true, uint16(len(wmName)),
			wmName).Reply()
		if err != nil {
			return "", err
		}

		response, err = xproto.GetProperty(X, false, win, wmNameAtom.Atom,
			xproto.GetPropertyTypeAny, 0, (1<<32)-1).Reply()
		if err != nil {
			return "", err
		}

		if response.Length == 0 {
			response = new(xproto.GetPropertyReply)
		}
	}
	return string(response.Value), nil
}

// Get the title of the currently focused window
func Title(X *xgb.Conn) (name string, err error) {

	// get the current window, and then use that
	// to get it's name

	currentId, err := GetFocus(X)
	if err != nil {
		return
	}

	name, err = GetName(X, currentId)
	if err != nil {
		return
	}

	// return the value of name
	return name, nil
}
