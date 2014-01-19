// Copyright 2013 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gettext

import (
	"fmt"
	"io/ioutil"
	"sync"
)

type domainManager struct {
	mutex        sync.Mutex
	locale       string
	domain       string
	domainPath   map[string]string
	domainData   map[string][]byte
	domainLocals map[string][]string
	trMap        map[string]*translator
}

func newDomainManager() *domainManager {
	return &domainManager{
		locale:       DefaultLocale,
		domainPath:   make(map[string]string),
		domainData:   make(map[string][]byte),
		domainLocals: make(map[string][]string),
		trMap:        make(map[string]*translator),
	}
}

func (p *domainManager) Bind(domain, path string, data []byte) (domains, paths []string, err error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	switch {
	case domain != "" && path != "": // bind new domain
		if _, ok := p.domainPath[domain]; ok {
			err = fmt.Errorf("gettext: domain already exists!")
			return
		}
		var locals []string
		if locals, err = p.globPathLocales(path); err != nil {
			return
		}
		for i := 0; i < len(locals); i++ {
			moName := p.getLocalFileName(domain, path, locals[i], ".mo")
			if f, err := newMoTranslator(moName, nil); err == nil { // ingore error
				key := p.makeTrMapKey(domain, locals[i])
				p.trMap[key] = f
				continue
			}
			poName := p.getLocalFileName(domain, path, locals[i], ".po")
			if f, err := newPoTranslator(poName, nil); err == nil { // ingore error
				key := p.makeTrMapKey(domain, locals[i])
				p.trMap[key] = f
				continue
			}
		}
		p.domainPath[domain] = path
		p.domainData[domain] = data
		p.domainLocals[domain] = locals
	case domain != "" && path == "": // delete domain
		if _, ok := p.domainPath[domain]; !ok {
			err = fmt.Errorf("gettext: domain not exists!")
			return
		}
		// enum locals
		var keys []string
		for _, v := range p.domainLocals[domain] {
			key := p.makeTrMapKey(domain, v)
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

func (p *domainManager) GetLocale() string {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.locale
}

func (p *domainManager) SetLocale(locale string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.locale = locale
	return nil
}

func (p *domainManager) GetDomain() string {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.domain
}

func (p *domainManager) SetDomain(domain string) error {
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

func (p *domainManager) Getdata(path string) []byte {
	panic("TODO")
}

func (p *domainManager) PNGettext(msgctxt, msgid, msgidPlural string, n int) string {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.gettext(p.domain, msgctxt, msgid, msgidPlural, n)
}

func (p *domainManager) DPNGettext(domain, msgctxt, msgid, msgidPlural string, n int) string {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.gettext(p.domain, msgctxt, msgid, msgidPlural, n)
}

func (p *domainManager) gettext(domain, msgctxt, msgid, msgidPlural string, n int) string {
	if p.locale == "" {
		return msgid
	}
	if f, ok := p.trMap[p.makeTrMapKey(domain, p.locale)]; ok {
		return f.PNGettext(msgctxt, msgid, msgidPlural, n)
	}
	return msgid
}

func (p *domainManager) globPathLocales(path string) (locals []string, err error) {
	list, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	for _, dir := range list {
		if dir.IsDir() {
			locals = append(locals, dir.Name())
		}
	}
	return
}

func (p *domainManager) getLocalFileName(domain, path, local, ext string) string {
	return fmt.Sprintf("%s/%s/LC_MESSAGES/%s%s", path, local, domain, ext)
}

func (p *domainManager) makeTrMapKey(domain, locale string) string {
	return domain + "_" + locale
}
