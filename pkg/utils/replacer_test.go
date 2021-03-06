// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"testing"
	"github.com/magiconair/properties/assert"
	"regexp"
	"os"
	"github.com/hidevopsio/hiboot/pkg/log"
	"reflect"
)

type Bar struct {
	Name    string
	Profile string
	SubBar  SubBar
	SubMap  map[string]interface{}
}

type Foo struct {
	Name    string
	Project string
	Bar     Bar
}

type SubBar struct {
	Name string
}

func TestParseReferences(t *testing.T) {
	s := []string{
		"name",
	}

	b := &SubBar{Name: "bar"}

	ref, err := ParseReferences(b, s)
	assert.Equal(t, nil, err)
	assert.Equal(t, "bar", ref)
}

func TestParseVariableName(t *testing.T) {

	re := regexp.MustCompile(`\$\{(.*?)\}`)
	matches := ParseVariables("foo ${name} bar", re)

	assert.Equal(t, "name", matches[0][1])
}

func TestGetFieldValue(t *testing.T) {
	foo := &Foo{Name: "foo"}

	field := GetFieldValue(foo, "Name")

	assert.Equal(t, reflect.String, field.Kind())
	assert.Equal(t, "foo", field.String())
}

func TestReplaceStringVariables(t *testing.T) {
	f := &Foo{
		Name: "foo",
		Bar: Bar{
			Name: "Hello ${name}",
		},
	}

	s, err := ReplaceStringVariables(f.Bar.Name, f)

	assert.Equal(t, nil, err)
	assert.Equal(t, "Hello foo", s)
}

func TestReplaceMap(t *testing.T) {
	b := &Bar{
		Name:    "bar",
		SubMap: map[string]interface{} {
			"name": "${name}",
			"nestedMap": map[string]interface{} {
				"name": "nested ${name}",
				"age": 18,
			},
		},
	}

	err := ReplaceMap(b.SubMap, b)
	assert.Equal(t, nil, err)
	assert.Equal(t, "bar", b.SubMap["name"])
	assert.Equal(t, "nested bar", b.SubMap["nestedMap"].(map[string]interface{})["name"])
}


func TestReplaceVariable(t *testing.T) {
	os.Setenv("FOO", "foo")
	os.Setenv("BAR", "bar")
	f := &Foo{
		Name:    "foo",
		Project: "it's ${FOO} project",
		Bar: Bar{
			Name:    "my name is ${BAR}",
			Profile: "${name}-bar",
			SubBar: SubBar{
				Name: "${bar.name}",
			},
			SubMap: map[string]interface{} {
				"barName": "${bar.name}",
				"name": "${name}",
				"nestedMap": map[string]interface{} {
					"name": "${name}",
					"age": 18,
				},
			},
		},
	}
	log.Println(f)
	err := Replace(f, f)
	log.Println(f)
	assert.Equal(t, nil, err)
	assert.Equal(t, "it's foo project", f.Project)
	assert.Equal(t, "foo-bar", f.Bar.Profile)
	assert.Equal(t, "my name is bar", f.Bar.Name)
	assert.Equal(t, f.Bar.Name, f.Bar.SubBar.Name)
	assert.Equal(t, f.Name, f.Bar.SubMap["name"])
	assert.Equal(t, f.Name, f.Bar.SubMap["nestedMap"].(map[string]interface{})["name"])
}

func TestParseVariables(t *testing.T) {
	os.Setenv("FOO", "foo")
	os.Setenv("BAR", "bar")
	os.Setenv("foo.bar", "fb")
	source := "the-${FOO}-${BAR}-${foo.bar}-env"

	re := regexp.MustCompile(`\$\{(.*?)\}`)

	matches := ParseVariables(source, re)

	log.Print(matches)
	assert.Equal(t, "${FOO}", matches[0][0])
	assert.Equal(t, "FOO", matches[0][1])
	assert.Equal(t, "${BAR}", matches[1][0])
	assert.Equal(t, "BAR", matches[1][1])
}

func TestReplaceReferences(t *testing.T) {
	os.Setenv("FOO", "foo")
	os.Setenv("BAR", "bar")
	f := &Foo{
		Name:    "foo",
		Project: "it's ${FOO} project",
		Bar: Bar{
			Name:    "my name is ${BAR}",
			Profile: "${name}-bar",
			SubBar: SubBar{
				Name: "${bar.name}",
			},
		},
	}
	res, err := ParseReferences(f, []string{"name"})
	assert.Equal(t, nil, err)
	log.Println("res: ", res)

	res, err = ParseReferences(f, []string{"bar", "name"})
	assert.Equal(t, nil, err)
	log.Println("res: ", res)
}
