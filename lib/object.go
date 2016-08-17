package lib

import (
	"fmt"
	"strings"
	"sort"
	"github.com/codahale/blake2"
)

const MaxObjectContent = 1024 * 1000

type Object struct {
	Name     string
	Metadata map[string]string
	MetadataAttr map[string]MetadataAttributes
	PrimaryKeyHash []byte
	Content  []byte
}

type MetadataAttributes struct {
	PrimaryKey bool
}

type LoadOptions map[string]interface{}

type SourceType interface {
	LoadObjects(source string, opt LoadOptions) ([]Object, error)
}

func (s *Seekret) LoadObjects(st SourceType, source string, opt LoadOptions) error {
	objectList, err := st.LoadObjects(source, opt)
	if err != nil {
		return err
	}
	s.objectList = append(s.objectList, objectList...)
	return nil
}

func (o *Object)SetMetadata(key string, value string, attr MetadataAttributes) error {
	if !validMetadataKey(key) {
		return fmt.Errorf("%s invalid key", key)
	}
	o.Metadata[key] = value
	o.MetadataAttr[key] = attr
	if attr.PrimaryKey {
		o.updatePrimaryKeyHash()
	}
	return nil
}

func (o *Object)GetMetadata(key string) (string,MetadataAttributes,error) {
	if !validMetadataKey(key) {
		return "",MetadataAttributes{},fmt.Errorf("%s invalid key", key)
	}

	val, ok := o.Metadata[key]
	if !ok {
		return "",MetadataAttributes{},fmt.Errorf("%s unexistent key", key)
	}

	attr, ok := o.MetadataAttr[key]
	if !ok {
		return "",MetadataAttributes{},fmt.Errorf("%s unexistent key", key)
	}

	return val,attr,nil
}

func (o *Object)GetMetadataAll(attr bool) (map[string]string) {
	return o.Metadata
}

func (o *Object)GetPrimaryKeyHash() []byte {
	return o.PrimaryKeyHash
}

func (o *Object)updatePrimaryKeyHash() {
	var primayKeyList []string
	for k,v := range o.MetadataAttr {
		if v.PrimaryKey {
			primayKeyList = append(primayKeyList, k)
		}
	}
	sort.Strings(primayKeyList)

	var text string
	for _,k := range primayKeyList {
		text = text + o.Metadata[k]
	}
	if text == "" {
		o.PrimaryKeyHash = nil
		return
	}

	h := blake2.New(&blake2.Config{
		Size: 32,
	})
	h.Write([]byte(text))
	o.PrimaryKeyHash = h.Sum(nil)
}

func validMetadataKey(key string) bool {
	return strings.Contains(key, ":")
}

/*
func GroupObjectsByMetadata(objects []Object, k string) (map[[]byte][]Object) {
	objectGroups := make(map[[]byte][]Object, 1, len(objects))
	for o := range objects {
		v := o.Metadata[k]
		objectList, ok = objectGroups[h]
	}
}

func GroupObjectsByPrimaryKeyHash(objects []Object) (map[[]byte][]Object) {
	objectGroups := make(map[[]byte][]Object, 1, len(objects))
	for o := range objects {
		h := o.PrimaryKeyHash(true)
		var objectList []Object
		objectList, ok = objectGroups[h]
		if !ok {
			objectList = make([]Object)
		}
	}
}
*/