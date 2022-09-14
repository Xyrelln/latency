// Copyright (C) 2015 The Android Open Source Project
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package adb

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

// func TestListDevices(t *testing.T) {
// 	devices, _ := Devices()
// 	fmt.Printf("list devices: %v", devices)
// 	// log.Printf("devices: %v", devices)

// }
func TestParseDevices(t *testing.T) {
	for i, c := range []struct {
		s string
		e error
		d []Device
	}{
		{
			d: []Device{{
				Serial: "02a2a2de20d7f6de",
				State:  Unauthorized,
			}},
			e: nil,
			s: `* daemon not running. starting it now on port 5037
* * daemon started successfully
* List of devices attached
02a2a2de20d7f6de        unauthorized
`,
		},
		{
			d: []Device{},
			e: nil,
			s: `List of devices attached
`,
		},
	} {
		devices, err := parseDevices(c.s)
		if c.e != err {
			t.Errorf("(%d) Expected error: %v, got error: %v", i, c.e, err)
		}
		count := len(devices)
		if count < len(c.d) {
			count = len(c.d)
		}
		for j := 0; j < count; j++ {
			var got, expected Device
			if j < len(devices) {
				got = *devices[j]
			}
			if j < len(c.d) {
				expected = c.d[j]
			}
			if got != expected {
				t.Errorf("(%d) Device %d was not as expected. Expected: %v, got: %v", i, j, expected, got)
			}
		}
	}
}
func (t treeNode) String() string {
	var s string
	if t.depth > 0 {
		s = strings.Repeat(" ", t.depth) + t.text
	} else {
		s = t.text
	}
	if cnt := len(t.children); cnt > 0 {
		c := make([]string, cnt)
		for i := range c {
			c[i] = t.children[i].String()
		}
		return s + "\n" + strings.Join(c, "\n")
	} else {
		return s
	}
}
func TestParseTabbedTree(t *testing.T) {
	expected := `
0
  00
    000
    001
  01
    010
    011
      0110
  02
1
  10`
	got := parseTabbedTree(expected).String()
	if got != expected {
		t.Errorf("Tree was not as expected.\nExpected: %v\nGot: %v", expected, got)
	}
}
func TestParseActions(t *testing.T) {
	str := `
Activity Resolver Table:
  Non-Data Actions:
    android.intent.action.MAIN:
      43178558 com.google.foo/.FooActivity filter 4327f110
    com.google.android.FOO:
      43178558 com.google.foo/.FooActivity filter 431d7db8
    android.intent.action.SEARCH:
      43178558 com.google.foo/.FooActivity filter 4327cc40
Packages:
  Package [com.google.foo] (ffffffc):
    userId=12345
    primaryCpuAbi=armeabi-v7a
    secondaryCpuAbi=null
    versionCode=902107 targetSdk=15
`
	expected := &InstalledPackage{
		Name: "com.google.foo",
		ABI:  "armeabi-v7a",
	}
	expected.Actions = []*Action{
		{
			Package:  expected,
			Name:     "android.intent.action.MAIN",
			Activity: ".FooActivity",
		}, {
			Package:  expected,
			Name:     "com.google.android.FOO",
			Activity: ".FooActivity",
		}, {
			Package:  expected,
			Name:     "android.intent.action.SEARCH",
			Activity: ".FooActivity",
		},
	}
	d := &Device{}
	packages, err := d.parsePackages(str)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(packages) != 1 {
		t.Errorf("Got %d packages, expected 1", len(packages))
	} else if len(packages[0].Actions) != len(expected.Actions) {
		t.Errorf("Got %d actions, expected %d", len(packages[0].Actions), len(expected.Actions))
	} else {
		for i, a := range packages[0].Actions {
			if a.Name != expected.Actions[i].Name {
				t.Errorf("[%d] Expected action %s got %s", i, a.Name, expected.Actions[i].Name)
			}
			if a.Activity != expected.Actions[i].Activity {
				t.Errorf("[%d] Expected activity %s got %s", i, a.Activity, expected.Actions[i].Activity)
			}
		}
		if packages[0].Name != expected.Name {
			t.Errorf("Expected package name %s got %s", expected.Name, packages[0].Name)
		}
		if packages[0].ABI != expected.ABI {
			t.Errorf("Expected ABI name %s got %s", expected.ABI, packages[0].ABI)
		}
	}
}

func TestParsePids(t *testing.T) {
	serial := "b9f8ef93"
	device := GetDevice(serial)
	device.KillScrcyServer()
}

func TestParseDisplay(t *testing.T) {
	serial := "b9f8ef93"
	device := GetDevice(serial)
	ds, _ := device.DisplaySize()
	fmt.Printf("current display siez:%v", ds)
}

func TestRegexp(t *testing.T) {
	// match, _ := regexp.MatchString("p([a-z]+)ch", "peach")
	// fmt.Println(match)

	// r, _ := regexp.Compile("p([a-z]+)ch")
	// fmt.Println(r.MatchString("peach"))

	re := regexp.MustCompile(`displayWidth=[0-9]+`)
	re2 := regexp.MustCompile(`displayHeight=[0-9]+`)
	fmt.Printf("%s \n", re.Find([]byte(`{displayWidth=1080 displayHeight=2400 density={3.0} cuto`)))
	fmt.Printf("%q\n", re2.FindAll([]byte(`{displayWidth=1080 displayHeight=2400 density={3.0} cuto`), -1))
}

func TestEnv(t *testing.T) {
	fmt.Println(os.Environ())
	fmt.Println(os.Getenv("ANDROID_HOME"))
	fmt.Println(exec.LookPath("adb"))

}
