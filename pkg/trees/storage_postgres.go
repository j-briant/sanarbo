package trees

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type PGX struct {
	con *pgxpool.Pool
	log *log.Logger
}

func (P PGX) List(offset, limit int) ([]*TreeList, error) {
	//TODO implement me
	panic("implement me")
}

func (P PGX) Get(id int32) (*Tree, error) {
	//TODO implement me
	panic("implement me")
}

func (P PGX) GetMaxId() (int32, error) {
	//TODO implement me
	panic("implement me")
}

func (P PGX) Exist(id int32) bool {
	//TODO implement me
	panic("implement me")
}

func (P PGX) Count() (int32, error) {
	//TODO implement me
	panic("implement me")
}

func (P PGX) Create(object Tree) (*Tree, error) {
	//TODO implement me
	panic("implement me")
}

func (P PGX) Update(id int32, object Tree) (*Tree, error) {
	//TODO implement me
	panic("implement me")
}

func (P PGX) Delete(id int32) error {
	//TODO implement me
	panic("implement me")
}

func (P PGX) SearchTreesByName(pattern string) ([]*TreeList, error) {
	//TODO implement me
	panic("implement me")
}

func (P PGX) IsTreeActive(id int32) bool {
	//TODO implement me
	panic("implement me")
}
