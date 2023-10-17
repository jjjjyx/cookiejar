package cookiejar

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"time"
)

func sortEntry(entries []entry) {
	sort.Slice(entries, func(i, j int) bool {
		e0, e1 := &entries[i], &entries[j]
		if e0.CanonicalHost != e1.CanonicalHost {
			return e0.CanonicalHost < e1.CanonicalHost
		}
		if len(e0.Path) != len(e1.Path) {
			return len(e0.Path) > len(e1.Path)
		}
		if !e0.Creation.Equal(e1.Creation) {
			return e0.Creation.Before(e1.Creation)
		}
		// The following are not strictly necessary
		// but are useful for providing deterministic
		// behaviour in tests.
		if e0.Name != e1.Name {
			return e0.Name < e1.Name
		}
		return e0.Value < e1.Value
	})
}

func entryToHttpCookie(entries []entry) []*http.Cookie {
	cookies := make([]*http.Cookie, len(entries))

	for i, e := range entries {
		cookie := &http.Cookie{
			Name:     e.Name,
			Value:    e.Value,
			Path:     e.Path,
			Domain:   e.Domain,
			Expires:  e.Expires,
			Secure:   e.Secure,
			HttpOnly: e.HttpOnly,
			MaxAge:   e.MaxAge,
		}
		switch e.SameSite {
		case "SameSite":
			cookie.SameSite = http.SameSiteDefaultMode
		case "SameSite=Strict":
			cookie.SameSite = http.SameSiteStrictMode
		case "SameSite=Lax":
			cookie.SameSite = http.SameSiteLaxMode
		}

		cookies[i] = cookie
	}

	return cookies
}

// GetAllCookies 获取所有添加进来的cookie, 保留 过期的cookie 和 因为maxAge 被删除的cookie
// 但是注意区别，如果本身就过期的cookie，或者maxAge <0 添加到jar 中就已经无法添加了 这里也无法导出这样的cookie
//
func (j *Jar) GetAllCookies() []*http.Cookie {
	var selected []entry
	j.mu.Lock()
	defer j.mu.Unlock()

	for _, submap := range j.entries {
		for _, e := range submap {
			selected = append(selected, e)
		}
	}
	sortEntry(selected)
	return entryToHttpCookie(selected)
}

// GetCookies 获取所有 cookie 但是忽略因为maxAge被删除 或者 因为超过过期时间的cookie
func (j *Jar) GetCookies() []*http.Cookie {
	now := time.Now()

	var selected []entry

	j.mu.Lock()
	defer j.mu.Unlock()
	for _, submap := range j.entries {
		for _, e := range submap {
			if !e.Expires.After(now) {
				// Do not return expired cookies.
				continue
			}
			selected = append(selected, e)
		}
	}

	sortEntry(selected)

	return entryToHttpCookie(selected)
}

// Save 序列化
func (j *Jar) Save() (string, error) {
	//entries := j.allPersistentEntries()
	ret, err := json.Marshal(j.entries)
	if err != nil {
		return "", err
	}

	return string(ret), err
}

// Load 导入序列化的结果
func (j *Jar) Load(data string) error {
	//var entries []entry
	if err := json.Unmarshal([]byte(data), &j.entries); err != nil {
		log.Printf("warning: discarding cookies in invalid format (error: %v)", err)
		return err
	}

	return nil
}
