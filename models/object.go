// Copyright 2016 - Authors included on AUTHORS file.
//
// Use of this source code is governed by a Apache License
// that can be found in the LICENSE file.

package models

import (
	"fmt"
	"github.com/codahale/blake2"
	"sort"
)

// MaxObjectContentLen contains the maximum size for the content of an object.
const MaxObjectContentLen = 1024 * 5000

// Contains a KeyHash or nil
type KeyHash *[]byte 

// Represents an object.
type Object struct {
	Name    string
	Content []byte

	Metadata       map[string]MetadataData
	PrimaryKeyHash KeyHash
}

// Represents the metadata of an object.
type MetadataData struct {
	value string
	attr  MetadataAttributes
}

// Represents the attributes of metadata.
type MetadataAttributes struct {
	// All objects with same value on this key has the same content. It's used
	// to optimise the inspection.
	PrimaryKey bool
}

// NewObject creates a new object.
func NewObject(name string, content []byte) *Object {
	if len(content) > MaxObjectContentLen {
		content = content[:MaxObjectContentLen]
	}
	o := &Object{
		Name:    name,
		Content: content,

		Metadata:       make(map[string]MetadataData),
		PrimaryKeyHash: nil,
	}
	return o
}

// SetMetadata sets a metadata value for the object.
func (o *Object) SetMetadata(key string, value string, attr MetadataAttributes) error {
	o.Metadata[key] = MetadataData{
		value: value,
		attr:  attr,
	}

	if attr.PrimaryKey {
		o.updatePrimaryKeyHash()
	}

	return nil
}

// SetMetadata gets a metadata value from the object.
func (o *Object) GetMetadata(key string) (string, error) {
	data, ok := o.Metadata[key]
	if !ok {
		return "", fmt.Errorf("%s unexistent key", key)
	}

	return data.value, nil
}

// GetMetadataAll gets a map that contains all metadata of the object.
func (o *Object) GetMetadataAll(attr bool) map[string]string {
	metadataAll := make(map[string]string)
	for k, v := range o.Metadata {
		metadataAll[k] = v.value
	}
	return metadataAll
}

// GetPrimaryKeyHash returns the primary key hash of the object. This hash is
// calculated by using the information of all metadata marked as primary key.
func (o *Object) GetPrimaryKeyHash() KeyHash {
	return o.PrimaryKeyHash
}

func (o *Object) updatePrimaryKeyHash() {
	var primayKeyList []string
	for k, v := range o.Metadata {
		if v.attr.PrimaryKey {
			primayKeyList = append(primayKeyList, k)
		}
	}
	sort.Strings(primayKeyList)

	var text string
	for _, k := range primayKeyList {
		text = text + fmt.Sprintf("{%s//%s}", k, o.Metadata[k].value)
	}
	if text == "" {
		o.PrimaryKeyHash = nil
		return
	}

	h := blake2.New(&blake2.Config{
		Size: 32,
	})
	h.Write([]byte(text))

	primayKeyHask := h.Sum(nil)
	o.PrimaryKeyHash = &primayKeyHask
}

func GroupObjectsByMetadata(objects []Object, k string) map[string][]Object {
	objectGroups := make(map[string][]Object)
	for _, o := range objects {
		v, err := o.GetMetadata(k)
		if err != nil {
			fmt.Println(err)
		}

		var objectList []Object
		var ok bool

		objectList, ok = objectGroups[v]
		if !ok {
			objectList = make([]Object, 0)
		}
		objectList = append(objectList, o)
		objectGroups[v] = objectList
	}
	return objectGroups
}

func GroupObjectsByPrimaryKeyHash(objects []Object) map[string][]Object {
	objectGroups := make(map[string][]Object)
	for _, o := range objects {
		var objectList []Object
		var ok bool
		var key string

		if o.PrimaryKeyHash != nil {
			key = string(*o.PrimaryKeyHash)
		} else {
			key = o.Name
		}

		objectList, ok = objectGroups[key]
		if !ok {
			objectList = make([]Object, 0)
		}
		objectList = append(objectList, o)
		objectGroups[key] = objectList
	}
	return objectGroups
}
