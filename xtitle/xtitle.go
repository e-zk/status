package xtitle

import (
	//"fmt"
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	//"github.com/BurntSushi/xgbutil"
	//"github.com/BurntSushi/xgbutil/ewmh"
	//"github.com/BurntSushi/xgbutil/xprop"
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

func GetFocus(X *xgb.Conn) xproto.Window {
	setup := xproto.Setup(X)
	rootId := setup.DefaultScreen(X).Root

	aname := "_NET_ACTIVE_WINDOW"
	activeAtom, err := xproto.InternAtom(X, true, uint16(len(aname)),
		aname).Reply()
	if err != nil {
		panic(err)
	}

	reply, err := xproto.GetProperty(X, false, rootId, activeAtom.Atom,
		xproto.GetPropertyTypeAny, 0, (1<<32)-1).Reply()
	if err != nil {
		panic(err)
	}

	currentId := xproto.Window(xgb.Get32(reply.Value))
	return currentId
}

func getName(X *xgb.Conn, win xproto.Window) (response *xproto.GetPropertyReply, err error) {

	err = nil
	// get name
	netWmName := "_NET_WM_NAME"
	netWmAtom, err := xproto.InternAtom(X, true, uint16(len(netWmName)),
		netWmName).Reply()
	if err != nil {
		// TODO error handling
		panic(err)
	}

	response, err = xproto.GetProperty(X, false, win, netWmAtom.Atom,
		xproto.GetPropertyTypeAny, 0, (1<<32)-1).Reply()
	if err != nil {
		panic(err)
	}

	// if there is no _NET_WM_NAME, attempt to get WM_NAME instead...
	if response.Length == 0 {
		wmName := "WM_NAME"
		wmNameAtom, err := xproto.InternAtom(X, true, uint16(len(wmName)),
			wmName).Reply()
		if err != nil {
			panic(err)
		}

		response, err = xproto.GetProperty(X, false, win, wmNameAtom.Atom,
			xproto.GetPropertyTypeAny, 0, (1<<32)-1).Reply()
		if err != nil {
			panic(err)
		}

		if response.Length == 0 {
			response = new(xproto.GetPropertyReply)
		}
	}
	return
}

func Title(X *xgb.Conn) string {
	// initialise connection to X
	/*X, err := xgb.NewConn()
	if err != nil {
		fmt.Printf("Error| %s\n", err)
	}

	// initialise Xutil
	Xutil, err := xgbutil.NewConnXgb(X)
	if err != nil {
		fmt.Printf("Error| %s\n", err)
	}*/

	// get the current window, and then use that
	// to get it's name
	currentId := GetFocus(X)
	name := getName(X, currentId).Value

	return string(name)
}

/*
func main() {
	// initialise connection to X
	X, err := xgb.NewConn()
	if err != nil {
		fmt.Printf("Error| %s\n", err)
	}

	// initialise Xutil
	Xutil, err := xgbutil.NewConnXgb(X)
	if err != nil {
		fmt.Printf("Error| %s\n", err)
	}

	// scan in wid from the last arg
	//fmt.Sscanf(os.Args[len(os.Args)-1], "0x%x", &wid)

	//
	currentId := getFocus(X)
	name := getName(Xutil, currentId).Value

	//
	fmt.Printf("%s\n", name)
}*/
