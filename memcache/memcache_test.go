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
	"testing"

	"github.com/go-chi/chi"
	. "github.com/smartystreets/goconvey/convey"
	"go.wandrs.dev/session"
)

func Test_MemcacheProvider(t *testing.T) {
	Convey("Test memcache session provider", t, func() {
		opt := session.Options{
			Provider:       "memcache",
			ProviderConfig: "127.0.0.1:9090",
		}

		Convey("Basic operation", func() {
			c := chi.NewRouter()
			c.Use(session.Sessioner(opt))

			c.Get("/", func(resp http.ResponseWriter, req *http.Request) {
				sess := session.GetSession(req)
				sess.Set("uname", "unknwon")
			})
			c.Get("/reg", func(resp http.ResponseWriter, req *http.Request) {
				sess := session.GetSession(req)
				raw, err := sess.RegenerateID(resp, req)
				So(err, ShouldBeNil)
				So(raw, ShouldNotBeNil)

				uname := raw.Get("uname")
				So(uname, ShouldNotBeNil)
				So(uname, ShouldEqual, "unknwon")
			})
			c.Get("/get", func(resp http.ResponseWriter, req *http.Request) {
				sess := session.GetSession(req)
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
			c.Use(session.Sessioner(opt))
			c.Get("/", func(resp http.ResponseWriter, req *http.Request) {
				sess := session.GetSession(req)
				raw, err := sess.RegenerateID(resp, req)
				So(err, ShouldBeNil)
				So(raw, ShouldNotBeNil)
			})

			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/", nil)
			So(err, ShouldBeNil)
			req.Header.Set("Cookie", "MacaronSession=ad2c7e3cbecfcf486; Path=/;")
			c.ServeHTTP(resp, req)
		})
	})
}
