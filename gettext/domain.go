// Copyright 2013 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gettext

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
)

var dTable = newDomainTable()

type domainTable struct {
	mutex        sync.Mutex
	locale       string
	domain       string
	domainPath   map[string]string
	domainLocals map[string][]string
	trMap        map[string]*translator
}

func makeDomainFileKey(domain, locale string) string {
	return domain + "_" + locale
}

func newDomainTable() *domainTable {
	return &domainTable{
		locale:       DefaultLocale,
		domainPath:   make(map[string]string),
		domainLocals: make(map[string][]string),
		trMap:        make(map[string]*translator),
	}
}

func (p *domainTable) Bind(domain, path string) (domains, paths []string, err error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	switch {
	case domain != "" && path != "": // bind new domain
		if _, ok := p.domainPath[domain]; ok {
			err = fmt.Errorf("gettext: domain already exists!")
			return
		}
		var locals, files []string
		locals, files, err = p.globDomainLocales(domain, path)
		if err != nil {
			return
		}
		for i := 0; i < len(files); i++ {
			if f, err := newMoTranslator(files[i]); err == nil { // ingore error
				key := makeDomainFileKey(domain, locals[i])
				p.trMap[key] = f
			}
		}
		p.domainPath[domain] = path
		p.domainLocals[domain] = locals
	case domain != "" && path == "": // delete domain
		if _, ok := p.domainPath[domain]; !ok {
			err = fmt.Errorf("gettext: domain not exists!")
			return
		}
		// enum locals
		var keys []string
		for _, v := range p.domainLocals[domain] {
			key := makeDomainFileKey(domain, v)
			keys = append(keys, key)
		}
		// delete all mo files
		for _, k := range keys {
			delete(p.trMap, k)
		}
		delete(p.domainLocals, domain)
		delete(p.domainPath, domain)
	}

	// return all bind domain
	for k, v := range p.domainPath {
		domains = append(domains, k)
		paths = append(paths, v)
	}
	return
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
	if domain != "" {
		if _, ok := p.domainPath[domain]; !ok {
			return fmt.Errorf("gettext: domain not exists!")
		}
	}
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
	if p.locale == "" {
		return msgid
	}
	if f, ok := p.trMap[makeDomainFileKey(domain, p.locale)]; ok {
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
		local := filepath.ToSlash(files[i])
		local = local[:strings.Index(local, "/LC_MESSAGES/"+domain+".mo")]
		local = local[strings.LastIndex(local, "/")+1:]
		locals = append(locals, local)
	}
	return
}
