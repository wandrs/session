// Copyright 2014 The Macaron Authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package session

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_Sessioner(t *testing.T) {
	Convey("Use session middleware", t, func() {
		c := chi.NewRouter()
		c.Use(Sessioner())
		c.Get("/", func(resp http.ResponseWriter, req *http.Request) {})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/", nil)
		So(err, ShouldBeNil)
		c.ServeHTTP(resp, req)
	})

	Convey("Register invalid provider", t, func() {
		Convey("Provider not exists", func() {
			defer func() {
				So(recover(), ShouldNotBeNil)
			}()

			c := chi.NewRouter()
			c.Use(Sessioner(Options{
				Provider: "fake",
			}))
		})

		Convey("Provider value is nil", func() {
			defer func() {
				So(recover(), ShouldNotBeNil)
			}()

			Register("fake", nil)
		})

		Convey("Register twice", func() {
			defer func() {
				So(recover(), ShouldNotBeNil)
			}()

			Register("memory", &MemProvider{})
		})
	})
}

func testProvider(opt Options) {
	Convey("Basic operation", func() {
		c := chi.NewRouter()
		c.Use(Sessioner(opt))

		c.Get("/", func(resp http.ResponseWriter, req *http.Request) {
			sess := GetSession(req)
			sess.Set("uname", "unknwon")
		})
		c.Get("/reg", func(resp http.ResponseWriter, req *http.Request) {
			sess := GetSession(req)
			raw, err := sess.RegenerateID(resp, req)
			So(err, ShouldBeNil)
			So(raw, ShouldNotBeNil)

			uname := raw.Get("uname")
			So(uname, ShouldNotBeNil)
			So(uname, ShouldEqual, "unknwon")
		})
		c.Get("/get", func(resp http.ResponseWriter, req *http.Request) {
			sess := GetSession(req)
			sid := sess.ID()
			So(sid, ShouldNotBeEmpty)

			raw, err := sess.Read(sid)
			So(err, ShouldBeNil)
			So(raw, ShouldNotBeNil)

			uname := sess.Get("uname")
			So(uname, ShouldNotBeNil)
			So(uname, ShouldEqual, "unknwon")

			So(sess.Delete("uname"), ShouldBeNil)
			So(sess.Get("uname"), ShouldBeNil)

			So(sess.Destroy(resp, req), ShouldBeNil)
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/", nil)
		So(err, ShouldBeNil)
		c.ServeHTTP(resp, req)

		cookie := resp.Header().Get("Set-Cookie")

		resp = httptest.NewRecorder()
		req, err = http.NewRequest("GET", "/reg", nil)
		So(err, ShouldBeNil)
		req.Header.Set("Cookie", cookie)
		c.ServeHTTP(resp, req)

		cookie = resp.Header().Get("Set-Cookie")

		resp = httptest.NewRecorder()
		req, err = http.NewRequest("GET", "/get", nil)
		So(err, ShouldBeNil)
		req.Header.Set("Cookie", cookie)
		c.ServeHTTP(resp, req)
	})

	Convey("Regenrate empty session", func() {
		c := chi.NewRouter()
		c.Use(Sessioner(opt))
		c.Get("/", func(resp http.ResponseWriter, req *http.Request) {
			sess := GetSession(req)
			raw, err := sess.RegenerateID(resp, req)
			So(err, ShouldBeNil)
			So(raw, ShouldNotBeNil)
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/", nil)
		So(err, ShouldBeNil)
		req.Header.Set("Cookie", "MacaronSession=ad2c7e3cbecfcf48; Path=/;")
		c.ServeHTTP(resp, req)
	})

	Convey("GC session", func() {
		c := chi.NewRouter()
		opt2 := opt
		opt2.Gclifetime = 1
		c.Use(Sessioner(opt2))

		c.Get("/", func(resp http.ResponseWriter, req *http.Request) {
			sess := GetSession(req)
			sess.Set("uname", "unknwon")
			So(sess.ID(), ShouldNotBeEmpty)
			uname := sess.Get("uname")
			So(uname, ShouldNotBeNil)
			So(uname, ShouldEqual, "unknwon")

			So(sess.Flush(), ShouldBeNil)
			So(sess.Get("uname"), ShouldBeNil)

			time.Sleep(2 * time.Second)
			sess.GC()
			So(sess.Count(), ShouldEqual, 0)
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/", nil)
		So(err, ShouldBeNil)
		c.ServeHTTP(resp, req)
	})

	Convey("Detect invalid sid", func() {
		c := chi.NewRouter()
		c.Use(Sessioner(opt))
		c.Get("/", func(resp http.ResponseWriter, req *http.Request) {
			sess := GetSession(req)
			raw, err := sess.Read("../session/ad2c7e3cbecfcf486")
			So(strings.Contains(err.Error(), "invalid 'sid': "), ShouldBeTrue)
			So(raw, ShouldBeNil)
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/", nil)
		So(err, ShouldBeNil)
		c.ServeHTTP(resp, req)
	})
}
