// Copyright 2013 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gettext

var (
	DefaultLocale = getDefaultLocale() // initialized with $(LC_MESSAGES) or $(LANG)
)

// SetLocale sets or queries the program's current locale.
//
// Examples:
//	SetLocale("")      // get locale: DefaultLocale
//	SetLocale("zh_CN") // set locale: return zh_CN
//	SetLocale("")      // get locale: return zh_CN
func SetLocale(locale string) (string, error) {
	if locale != "" {
		if err := dTable.SetLocale(locale); err != nil {
			return "", err
		}
	}
	return dTable.GetLocale(), nil
}

// Textdomain sets or retrieves the current message domain.
//
// Examples:
//	Textdomain("poedit") // set domain
//	Textdomain("")       // get domain
func Textdomain(domain string) (string, error) {
	if domain != "" {
		if err := dTable.SetDomain(domain); err != nil {
			return "", err
		}
	}
	return dTable.GetDomain(), nil
}

// BindTextdomain sets directory containing message.
//
// Examples:
//	BindTextdomain("poedit", "local")
func BindTextdomain(domain, path string) error {
	return dTable.Bind(domain, path)
}

// Gettext attempt to translate a text string into the user's native language,
// by looking up the translation in a message catalog.
//
// It use the caller's function name as the msgctxt.
//
// Examples:
//	func Foo() {
//		msg := gettext.Gettext("Hello") // msgctxt is "some/package/name.Foo"
//	}
func Gettext(msgid string) string {
	return PGettext(callerName(2), msgid)
}

// NGettext attempt to translate a text string into the user's native language,
// by looking up the appropriate plural form of the translation in a message
// catalog.
//
// It use the caller's function name as the msgctxt.
//
// Examples:
//	func Foo() {
//		msg := gettext.NGettext("%d people", "%d peoples", 2)
//	}
func NGettext(msgid, msgidPlural string, n int) string {
	return PNGettext(callerName(2), msgid, msgidPlural, n)
}

// Gettext attempt to translate a text string into the user's native language,
// by looking up the translation in a message catalog.
//
// Examples:
//	func Foo() {
//		msg := gettext.PGettext("gettext-go.example", "Hello") // msgctxt is "gettext-go.example"
//	}
func PGettext(msgctxt, msgid string) string {
	return PNGettext(msgctxt, msgid, "", 0)
}

// PNGettext attempt to translate a text string into the user's native language,
// by looking up the appropriate plural form of the translation in a message
// catalog.
//
// Examples:
//	func Foo() {
//		msg := gettext.PNGettext("gettext-go.example", "%d people", "%d peoples", 2)
//	}
func PNGettext(msgctxt, msgid, msgidPlural string, n int) string {
	return dTable.PNGettext(msgctxt, msgid, msgidPlural, n)
}

// DGettext like Gettext(), but looking up the message in the specified domain.
//
// Examples:
//	func Foo() {
//		msg := gettext.DGettext("poedit", "Hello")
//	}
func DGettext(domain, msgid string) string {
	return DPGettext(domain, callerName(2), msgid)
}

// DGettext like NGettext(), but looking up the message in the specified domain.
//
// Examples:
//	func Foo() {
//		msg := gettext.PNGettext("poedit", "gettext-go.example", "%d people", "%d peoples", 2)
//	}
func DNGettext(domain, msgid, msgidPlural string, n int) string {
	return DPNGettext(domain, callerName(2), msgid, msgidPlural, n)
}

// DGettext like PGettext(), but looking up the message in the specified domain.
//
// Examples:
//	func Foo() {
//		msg := gettext.DPGettext("poedit", "gettext-go.example", "Hello")
//	}
func DPGettext(domain, msgctxt, msgid string) string {
	return DPNGettext(domain, callerName(2), msgid, "", 0)
}

// DPNGettext like PNGettext(), but looking up the message in the specified domain.
//
// Examples:
//	func Foo() {
//		msg := gettext.DPNGettext("poedit", "gettext-go.example", "%d people", "%d peoples", 2)
//	}
func DPNGettext(domain, msgctxt, msgid, msgidPlural string, n int) string {
	return dTable.DPNGettext(domain, msgctxt, msgid, msgidPlural, n)
}
