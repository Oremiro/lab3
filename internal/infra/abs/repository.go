package abs

type Repository interface {
	Insert(interface{}) error
	FindOneByField(string, interface{}, interface{}) error
}
