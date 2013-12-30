// Copyright 2013 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gettext

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"code.google.com/p/gettext-go/gettext/mo"
)

var dTable = newDomainTable()

type domainTable struct {
	mutex           sync.Mutex
	locale          string
	domain          string
	domainPath      map[string]string
	domainLocalFile map[string]*mo.File
}

func makeDomainFileKey(domain, locale string) string {
	return domain + "_" + locale
}

func newDomainTable() *domainTable {
	return &domainTable{
		locale:          DefaultLocale,
		domainPath:      make(map[string]string),
		domainLocalFile: make(map[string]*mo.File),
	}
}

func (p *domainTable) Bind(domain, path string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.domainPath[domain]; ok {
		return fmt.Errorf("gettext: domain already exists!")
	}
	locals, files, err := p.globDomainLocales(domain, path)
	if err != nil {
		return err
	}
	for i := 0; i < len(files); i++ {
		if f, err := mo.Load(files[i], nil); err == nil { // ingore error
			key := makeDomainFileKey(domain, locals[i])
			p.domainLocalFile[key] = f
		}
	}
	p.domainPath[domain] = path
	return nil
}

func (p *domainTable) GetLocale() string {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.locale
}

func (p *domainTable) SetLocale(locale string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.locale = locale
	return nil
}

func (p *domainTable) GetDomain() string {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.domain
}

func (p *domainTable) SetDomain(domain string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.domain = domain
	return nil
}

func (p *domainTable) PNGettext(msgctxt, msgid, msgidPlural string, n int) string {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.gettext(p.domain, msgctxt, msgid, msgidPlural, n)
}

func (p *domainTable) DPNGettext(domain, msgctxt, msgid, msgidPlural string, n int) string {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.gettext(p.domain, msgctxt, msgid, msgidPlural, n)
}

func (p *domainTable) gettext(domain, msgctxt, msgid, msgidPlural string, n int) string {
	if f, ok := p.domainLocalFile[makeDomainFileKey(domain, p.locale)]; ok {
		return f.PNGettext(msgctxt, msgid, msgidPlural, n)
	}
	return msgid
}

func (p *domainTable) globDomainLocales(domain, path string) (locals, files []string, err error) {
	pattern := filepath.Join(path, "*", "LC_MESSAGES", domain+".mo")
	if files, err = filepath.Glob(pattern); err != nil {
		return
	}
	for i := 0; i < len(files); i++ {
		local := files[i]
		local = local[:strings.Index(local, "/LC_MESSAGES/"+domain+".mo")]
		local = local[strings.LastIndex(local, "/")+1:]
		locals = append(locals, local)
	}
	return
}
