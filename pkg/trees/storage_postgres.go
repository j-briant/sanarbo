package trees

import (
	"context"
	"errors"
	"log"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PGX struct {
	con *pgxpool.Pool
	log *log.Logger
}

func (P PGX) List(offset, limit int) ([]*TreeList, error) {
	P.log.Printf("trace: entering List(%d, %d)", offset, limit)
	var res []*TreeList

	err := pgxscan.Select(context.Background(), P.con, &res, treesList, limit, offset)
	if err != nil {
		P.log.Printf("error: List pgxscan.Select unexpectedly failed, error : %v", err)
		return nil, err
	}
	if res == nil {
		P.log.Println("info : List returned no results ")
		return nil, errors.New("records not found")
	}

	return res, nil
}

func (P PGX) Get(id int32) (*Tree, error) {
	P.log.Printf("trace : entering Get(%d)", id)
	res := &Tree{}
	
	err := pgxscan.Get(context.Background(), P.con, res, treesGet, id)
	if err != nil {
		P.log.Printf("error : Get(%d) pgxscan.Select unexpectedly failed, error : %v", id, err)
		return nil, err
	}
	if res == nil {
		P.log.Printf("info : Get(%d) returned no results ", id)
		return nil, errors.New("records not found")
	}
	return res, nil
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
	P.log.Printf("trace : entering Create(%q,%q,%#v)", object.Name, object.Geom, object.TreeAttributes)
	var lastInsertId int = 0

	err := P.con.QueryRow(context.Background(), treesCreate, 
		object.Name, &object.Description, object.ExternalId, object.IsActive, &object.Comment, object.Creator, object.Geom, object.TreeAttributes).Scan(&lastInsertId)
	if err != nil {
		P.log.Printf("error : Create(%q) unexpectedly failed. error : %v", object.Name, err)
		return nil, err
	}
	P.log.Printf("info : Create(%q) created with id : %v", object.Name, lastInsertId)

	createdTree, err := P.Get(int32(lastInsertId))
	if err != nil {
		return nil, GetErrorF("error : tree was created, but cannot be retrieved", err)
	}
	return createdTree, nil
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

func (P PGX) IsUserAdmin(id int32) bool {
	//TODO implement a better user check...
	//Now only user with id(999) (bill board) is considered as admin
	if id == 999 {
		return true
	} else {
		return false
	}
}