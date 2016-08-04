package lib

const MaxObjectContent = 1024 * 1000

type Object struct {
	Name     string
	Metadata map[string]string
	Content  []byte
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