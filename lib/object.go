package lib

import (
	"fmt"
	"sort"
	"github.com/codahale/blake2"
)

const MaxObjectContentLen = 1024 * 5000

type Object struct {
	Name     string
	Content  []byte

	Metadata map[string]MetadataData
	PrimaryKeyHash []byte
}

type MetadataData struct {
	value string
	attr MetadataAttributes
}

type MetadataAttributes struct {
	PrimaryKey bool
}

func NewObject(name string, content []byte) *Object {
	if len(content) > MaxObjectContentLen {
		content = content[:MaxObjectContentLen]
	}
	o := &Object{
		Name: name,
		Content: content,

		Metadata: make(map[string]MetadataData),
		PrimaryKeyHash: nil,
	}
	return o
}

func (o *Object)SetMetadata(key string, value string, attr MetadataAttributes) error {
	o.Metadata[key] = MetadataData{
		value: value,
		attr: attr,
	}

	if attr.PrimaryKey {
		o.updatePrimaryKeyHash()
	}

	return nil
}

func (o *Object)GetMetadata(key string) (string,error) {
	data, ok := o.Metadata[key]
	if !ok {
		return "",fmt.Errorf("%s unexistent key", key)
	}

	return data.value,nil 
}

func (o *Object)GetMetadataAll(attr bool) (map[string]string) {
	metadataAll := make(map[string]string)
	for k,v := range o.Metadata {
		metadataAll[k] = v.value
	}
	return metadataAll
}

func (o *Object)GetPrimaryKeyHash() []byte {
	return o.PrimaryKeyHash
}

func (o *Object)updatePrimaryKeyHash() {
	var primayKeyList []string
	for k,v := range o.Metadata{
		if v.attr.PrimaryKey {
			primayKeyList = append(primayKeyList, k)
		}
	}
	sort.Strings(primayKeyList)

	var text string
	for _,k := range primayKeyList {
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
	o.PrimaryKeyHash = h.Sum(nil)
}


func GroupObjectsByMetadata(objects []Object, k string) (map[string][]Object) {
	objectGroups := make(map[string][]Object)
	for _,o := range objects {
		v,err := o.GetMetadata(k)
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



func GroupObjectsByPrimaryKeyHash(objects []Object) (map[string][]Object) {
	objectGroups := make(map[string][]Object)
	for _,o := range objects {
		var objectList []Object
		var ok bool

		objectList, ok = objectGroups[string(o.PrimaryKeyHash)]
		if !ok {
			objectList = make([]Object, 0)
		}
		objectList = append(objectList, o)
		objectGroups[string(o.PrimaryKeyHash)] = objectList
	}
	return objectGroups
}

