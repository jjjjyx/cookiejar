package cookiejar

import (
	"fmt"
	"net/http"
	"time"

	"net/url"
	"testing"
)

func TestGetAllCookies(t *testing.T) {
	jar, err := New(nil)
	if err != nil {
		fmt.Println("Error creating cookie jar:", err)
		return
	}

	// Define a sample URL
	baseURL, _ := url.Parse("https://example.com")

	// Add some cookies to the jar (you can add your own cookies)
	cookie1 := &http.Cookie{
		Name:   "cookie1",
		Value:  "value1",
		Domain: "example.com",
		//Expires:  someTime, // Set the expiration time
		HttpOnly: true,
		Secure:   true,
		Path:     "/path",
	}
	// 过期的
	cookie2 := &http.Cookie{
		Name:    "cookie2",
		Value:   "value2",
		Domain:  "example.com",
		Expires: time.Now().Add(1 * time.Second),
		//Expires:  someTime, // Set the expiration time
		HttpOnly: true,
		Secure:   true,
		Path:     "/path",
	}

	// 跨域的域名
	cookie3 := &http.Cookie{
		Name:   "cookie3",
		Value:  "value3",
		Domain: "www.abc.com",
		//Expires:  someTime, // Set the expiration time
		HttpOnly: true,
		Secure:   true,
		Path:     "/path",
	}

	//// maxAge
	//cookie3 := &http.Cookie{
	//	Name:   "cookie3",
	//	Value:  "value3",
	//	Domain: "www.abc.com",
	//	//Expires:  someTime, // Set the expiration time
	//	HttpOnly: true,
	//	Secure:   true,
	//	Path:     "/path",
	//}

	jar.SetCookies(baseURL, []*http.Cookie{cookie1, cookie2, cookie3})

	cookies := jar.GetAllCookies()

	if len(cookies) != 2 {
		t.Errorf("cookie 的数量应该是 2, 实际 = %d", len(cookies))
	}

	time.Sleep(2 * time.Second)

	cookies = jar.GetAllCookies()

	if len(cookies) != 2 {
		t.Errorf("cookie 的数量应该是 2, 实际 = %d", len(cookies))
	}

	cookies = jar.GetCookies()

	if len(cookies) != 1 {
		t.Errorf("cookie 的数量应该是 1, 实际 = %d", len(cookies))
	}
	//fmt.Println(cookies)

}

func TestMaxAge(t *testing.T) {
	jar, err := New(nil)
	if err != nil {
		fmt.Println("Error creating cookie jar:", err)
		return
	}

	// Define a sample URL
	baseURL, _ := url.Parse("https://example.com")

	// Add some cookies to the jar (you can add your own cookies)
	cookie1 := &http.Cookie{
		Name:   "cookie1",
		Value:  "value1",
		Domain: "example.com",

		//Expires:  someTime, // Set the expiration time
		HttpOnly: true,
		Secure:   true,
		Path:     "/path",
	}
	cookie2 := &http.Cookie{
		Name:   "cookie2",
		Value:  "value2",
		Domain: "example.com",
		MaxAge: -1,
		//Expires:  someTime, // Set the expiration time
		HttpOnly: true,
		Secure:   true,
		Path:     "/path",
	}

	cookie3 := &http.Cookie{
		Name:   "cookie3",
		Value:  "value3",
		Domain: "example.com",
		MaxAge: 2, // 2秒后过期
		//Expires:  someTime, // Set the expiration time
		HttpOnly: true,
		Secure:   true,
		Path:     "/path",
	}

	jar.SetCookies(baseURL, []*http.Cookie{cookie1, cookie2, cookie3})

	cookies := jar.GetAllCookies()

	if len(cookies) != 2 {
		t.Errorf("cookie 的数量应该是 2")
	}

	time.Sleep(2 * time.Second)

	cookies = jar.GetAllCookies()

	if len(cookies) != 2 {
		t.Errorf("cookie 的数量应该是 2, 实际 = %d", len(cookies))
	}

	cookies = jar.GetCookies()

	if len(cookies) != 1 {
		t.Errorf("cookie 的数量应该是 1, 实际 = %d", len(cookies))
	}
	//fmt.Println(cookies)

}

func TestSave(t *testing.T) {
	jar, err := New(nil)

	if err != nil {
		fmt.Println("Error creating cookie jar:", err)
		return
	}

	// Define a sample URL
	baseURL, _ := url.Parse("https://example.com")

	// Add some cookies to the jar (you can add your own cookies)
	cookie1 := &http.Cookie{
		Name:   "cookie1",
		Value:  "value1",
		Domain: "example.com",

		//Expires:  someTime, // Set the expiration time
		HttpOnly: true,
		Secure:   true,
		Path:     "/path",
	}
	cookie2 := &http.Cookie{
		Name:   "cookie2",
		Value:  "value2",
		Domain: "example.com",
		MaxAge: -1,
		//Expires:  someTime, // Set the expiration time
		HttpOnly: true,
		Secure:   true,
		Path:     "/path",
	}

	cookie3 := &http.Cookie{
		Name:   "cookie3",
		Value:  "value3",
		Domain: "example.com",
		MaxAge: 2, // 2秒后过期
		//Expires:  someTime, // Set the expiration time
		HttpOnly: true,
		Secure:   true,
		Path:     "/path",
	}

	jar.SetCookies(baseURL, []*http.Cookie{cookie1, cookie2, cookie3})

	save, _ := jar.Save()
	fmt.Println(save)
	jar2, _ := New(nil)

	_ = jar2.Load(save)

	c1 := jar.GetCookies()
	c2 := jar2.GetCookies()

	if len(c1) != len(c2) {
		t.Errorf("导入后cookie 的数量应该是相同, 实际 = %d - %d", len(c1), len(c2))
	}
	for i := 0; i < len(c1); i++ {
		co1 := c1[i].String()
		co2 := c2[i].String()

		if co1 != co2 {
			t.Errorf("cookie 值应该相同, 实际 = %s ---- %s", co1, co2)
		}
	}

}

func TestDomain(t *testing.T) {
	jar, err := New(nil)
	if err != nil {
		fmt.Println("Error creating cookie jar:", err)
		return
	}

	// Define a sample URL
	baseURL, _ := url.Parse("https://xxxx.example.com")

	// Add some cookies to the jar (you can add your own cookies)
	cookie1 := &http.Cookie{
		Name:   "cookie1",
		Value:  "value1",
		Domain: ".example.com",

		//Expires:  someTime, // Set the expiration time
		HttpOnly: true,
		Secure:   true,
		Path:     "/path",
	}

	jar.SetCookies(baseURL, []*http.Cookie{cookie1})

	cookies := jar.GetAllCookies()

	if len(cookies) != 1 {
		t.Errorf("cookie 的数量应该是 1, 实际 = %d", len(cookies))
	}
	cookie := cookies[0]
	if cookie.Domain != ".example.com" {
		fmt.Println(cookie.Domain)
		t.Errorf("cookie 的域名应该是 '.example.com', 实际 = %s", cookie.Domain)
	}
	//fmt.Println(cookies)

}
