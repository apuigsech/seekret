// Copyright 2016 - Authors included on AUTHORS file.
//
// Use of this source code is governed by a Apache License
// that can be found in the LICENSE file.

package models

import (
	"bytes"
	"testing"
)

type NewObjectSample struct {
	name       string
	expName    string
	content    []byte
	expContent []byte
	ok         bool
}

func TestNewObject(t *testing.T) {
	testSamples := []NewObjectSample{
		{
			name:       "test_1",
			expName:    "test_1",
			content:    []byte("content_1"),
			expContent: []byte("content_1"),
			ok:         true,
		},
		{
			name:       "test_2",
			expName:    "xxx",
			content:    []byte("content_2"),
			expContent: []byte("content_2"),
			ok:         false,
		},
		{
			name:       "test_3",
			expName:    "test_3",
			content:    []byte("content_3"),
			expContent: []byte("xxx"),
			ok:         false,
		},
	}

	for _, ts := range testSamples {
		if !testNewObjectSample(ts) {
			t.Error("unexpected new object")
		}
	}
}

func testNewObjectSample(ts NewObjectSample) bool {
	o := NewObject(ts.name, ts.content)

	ok := (o.Name == ts.expName && bytes.Equal(o.Content, ts.expContent))

	if ok == ts.ok {
		return true
	} else {
		return false
	}
}

type MetadataSample struct {
	key        string
	value      string
	expValue   string
	primaryKey bool
	ok         bool
}

func TestMetadata(t *testing.T) {
	testSamples := []MetadataSample{
		{
			key:        "key_1",
			value:      "value_1",
			expValue:   "value_1",
			primaryKey: false,
			ok:         true,
		},
		{
			key:        "key_2",
			value:      "value_2",
			expValue:   "value_2",
			primaryKey: true,
			ok:         true,
		},
	}

	for _, ts := range testSamples {
		if !testMetadataSample(ts) {
			t.Error("unexpected metadata")
		}
	}

}

func testMetadataSample(ts MetadataSample) bool {
	o := NewObject("test", []byte("content"))

	o.SetMetadata(ts.key, ts.value, MetadataAttributes{PrimaryKey: ts.primaryKey})
	value, err := o.GetMetadata(ts.key)

	ok := (err == nil && value == ts.expValue)

	if ok == ts.ok {
		return true
	} else {
		return false
	}
}

type PrimaryKeyHashSample struct {
	metadata       map[string]MetadataData
	primaryKeyHash KeyHash
	ok             bool
}

func TestPrimaryKeyHash(t *testing.T) {
	testSamples := []PrimaryKeyHashSample{
		{
			metadata:       map[string]MetadataData{},
			primaryKeyHash: nil,
			ok:             true,
		},
		{
			metadata: map[string]MetadataData{
				"key_1": MetadataData{
					value: "value_1",
					attr:  MetadataAttributes{PrimaryKey: false},
				},
			},
			primaryKeyHash: nil,
			ok:             true,
		},
		{
			metadata: map[string]MetadataData{
				"key_1": MetadataData{
					value: "value_1",
					attr:  MetadataAttributes{PrimaryKey: true},
				},
			},
			primaryKeyHash: &[]byte{0xe5, 0xd7, 0x3d, 0xe2, 0x2b, 0xa, 0xb1, 0x2, 0x64, 0xbb, 0x9, 0x77, 0xae, 0xea, 0x7, 0x4f, 0xd7, 0x14, 0x5d, 0xeb, 0x93, 0x84, 0xdc, 0xe, 0xb0, 0x91, 0x37, 0x29, 0x10, 0x56, 0x3, 0x45},
			ok:             true,
		},
		{
			metadata: map[string]MetadataData{
				"key_1": MetadataData{
					value: "value_1",
					attr:  MetadataAttributes{PrimaryKey: true},
				},
				"key_2": MetadataData{
					value: "value_2",
					attr:  MetadataAttributes{PrimaryKey: false},
				},
				"key_3": MetadataData{
					value: "value_3",
					attr:  MetadataAttributes{PrimaryKey: true},
				},
			},
			primaryKeyHash: &[]byte{0xd1, 0x9c, 0xb2, 0x5e, 0x12, 0x64, 0x82, 0xfb, 0xc8, 0x96, 0xc4, 0x9a, 0x65, 0xd0, 0x4e, 0xa2, 0xb8, 0x9, 0x67, 0x41, 0xdc, 0xa1, 0xc4, 0x82, 0x3d, 0x6f, 0xa1, 0x1, 0x1d, 0xaa, 0x45, 0xbc},
			ok:             true,
		},
		{
			metadata: map[string]MetadataData{
				"key_3": MetadataData{
					value: "value_3",
					attr:  MetadataAttributes{PrimaryKey: true},
				},
				"key_2": MetadataData{
					value: "value_2",
					attr:  MetadataAttributes{PrimaryKey: false},
				},
				"key_1": MetadataData{
					value: "value_1",
					attr:  MetadataAttributes{PrimaryKey: true},
				},
			},
			primaryKeyHash: &[]byte{0xd1, 0x9c, 0xb2, 0x5e, 0x12, 0x64, 0x82, 0xfb, 0xc8, 0x96, 0xc4, 0x9a, 0x65, 0xd0, 0x4e, 0xa2, 0xb8, 0x9, 0x67, 0x41, 0xdc, 0xa1, 0xc4, 0x82, 0x3d, 0x6f, 0xa1, 0x1, 0x1d, 0xaa, 0x45, 0xbc},
			ok:             true,
		},
	}

	for _, ts := range testSamples {
		if !testPrimaryKeyHashSample(ts) {
			t.Error("unexpected primary key hash")
		}
	}

}

func testPrimaryKeyHashSample(ts PrimaryKeyHashSample) bool {
	o := NewObject("test", []byte("content"))

	for k, v := range ts.metadata {
		o.SetMetadata(k, v.value, v.attr)
	}

	h := o.GetPrimaryKeyHash()

	var ok bool
	if h == nil {
		if ts.primaryKeyHash == nil {
			ok = true
		} else {
			ok = false
		}
	} else {
		ok = (bytes.Equal(*ts.primaryKeyHash, *h))
	}

	if ok == ts.ok {
		return true
	} else {
		return false
	}
}
