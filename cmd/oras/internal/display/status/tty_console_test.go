//go:build freebsd || linux || netbsd || openbsd || solaris

/*
Copyright The ORAS Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package status

import (
	"testing"

	"oras.land/oras-go/v2/content/memory"
	"oras.land/oras/internal/testutils"
)

func TestTTYPushHandler_TrackTarget(t *testing.T) {
	// prepare pty
	_, slave, err := testutils.NewPty()
	if err != nil {
		t.Fatal(err)
	}
	defer slave.Close()
	ph := NewTTYPushHandler(slave, mockFetcher.Fetcher)
	store := memory.New()
	// test
	_, fn, err := ph.TrackTarget(store)
	if err != nil {
		t.Fatal("TrackTarget() should not return an error")
	}
	defer func() {
		if err := fn(); err != nil {
			t.Fatal(err)
		}
	}()
	if ttyPushHandler, ok := ph.(*TTYPushHandler); !ok {
		t.Fatalf("TrackTarget() should return a *TTYPushHandler, got %T", ttyPushHandler)
	}
}

func Test_TTYPullHandler_TrackTarget(t *testing.T) {
	src := memory.New()
	t.Run("has TTY", func(t *testing.T) {
		_, device, err := testutils.NewPty()
		if err != nil {
			t.Fatal(err)
		}
		defer device.Close()
		ph := NewTTYPullHandler(device)
		got, fn, err := ph.TrackTarget(src)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err := fn(); err != nil {
				t.Fatal(err)
			}
		}()
		if got == src {
			t.Fatal("GraphTarget not be modified on TTY")
		}
	})

	t.Run("invalid TTY", func(t *testing.T) {
		ph := NewTTYPullHandler(nil)

		if _, _, err := ph.TrackTarget(src); err == nil {
			t.Fatal("expected error for no tty but got nil")
		}
	})
}
