package sqlBackend

type sql struct {
}

func New(connectionString string) (error, sql) {
	//condb, errdb := SQL.Open("mssql", "server=localhost;user id=sa;password=SA_PASSWORD=yourStrong(!)Password;")
	return nil, sql{}
}

//match the interface
func (sqlStruct sql) Get(id string) string {
	panic("not implemented")
}

//match the interface
func (sqlStruct sql) Add(uRL string) string {
	panic("not implemented")
}
