package client

// DataSource a source of binary data
type DataSource interface {
	// ObjectPath a unique path representation for the data, can be used as 'object path' when uploading data to lcoud storage
	ObjectPath() *ObjectPath
	// Data binary data
	Data() []byte
	// AddError add error - e.g.if any problems processing binary data
	AddError(err error)
}

// DataSourceIterator iterate over a collection of DataSource's
type DataSourceIterator interface {
	Each(handler func(DataSource))
	Size() int
}

// TaskDataSourceIterator a DataSource iterator
type TaskDataSourceIterator []Task

// Each handle DataSource iteration
func (as TaskDataSourceIterator) Each(handler func(DataSource)) {
	for _, a := range as {
		handler(a)
	}
}

// Size number of DataSources
func (as TaskDataSourceIterator) Size() int {
	return len(as)
}
