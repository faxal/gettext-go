// Copyright 2013 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gettext

var (
	DefaultLocale = getDefaultLocale() // default local is $(LC_MESSAGES) or $(LANG)
)

func init() {
	if DefaultLocale != "" {
		SetLocale(DefaultLocale)
	}
}

// SetLocale sets and queries the program's current locale.
//
// If the locale is not empty string, set the new local.
//
// If the locale is empty string, don't change anything.
//
// Returns is the current locale.
//
// Examples:
//	SetLocale("")      // get locale: return DefaultLocale
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

// BindTextdomain sets and queries program's domains.
//
// If the domain and path are all not empty string, bind the new domain.
// If the domain already exists, return error.
//
// If the domain is not empty string, but the path is the empty string,
// delete the domain.
// If the domain don't exists, return error.
//
// If the domain and the path are all empty string, don't change anything.
//
// Returns is the all bind domains.
//
// Examples:
//	BindTextdomain("poedit", "local") // bind "poedit" domain
//	BindTextdomain("", "")            // return all domains
//	BindTextdomain("poedit", "")      // delete "poedit" domain
//	BindTextdomain("", "")            // return all domains
func BindTextdomain(domain, path string) (domains, paths []string, err error) {
	return dTable.Bind(domain, path)
}

// Textdomain sets and retrieves the current message domain.
//
// If the domain is not empty string, set the new domains.
//
// If the domain is empty string, don't change anything.
//
// Returns is the all used domains.
//
// Examples:
//	Textdomain("poedit") // set domain: poedit
//	Textdomain("")       // get domain: return poedit
func Textdomain(domain string) (string, error) {
	if domain != "" {
		if err := dTable.SetDomain(domain); err != nil {
			return "", err
		}
	}
	return dTable.GetDomain(), nil
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

// PGettext attempt to translate a text string into the user's native language,
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

// DNGettext like NGettext(), but looking up the message in the specified domain.
//
// Examples:
//	func Foo() {
//		msg := gettext.PNGettext("poedit", "gettext-go.example", "%d people", "%d peoples", 2)
//	}
func DNGettext(domain, msgid, msgidPlural string, n int) string {
	return DPNGettext(domain, callerName(2), msgid, msgidPlural, n)
}

// DPGettext like PGettext(), but looking up the message in the specified domain.
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
