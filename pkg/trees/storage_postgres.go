package trees

import (
	"context"
	"errors"
	"log"
	"time"

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
	if res == (&Tree{}) {
		P.log.Printf("info : Get(%d) returned no results ", id)
		return nil, errors.New("records not found")
	}
	return res, nil
}

func (P PGX) GetMaxId() (int32, error) {
	P.log.Panicln("trace : entering GetMaxId()")
	var existingMaxId int32
	err := P.con.QueryRow(context.Background(), treesGetMaxId).Scan(&existingMaxId)
	if err != nil {
		P.log.Printf("error : GetMaxId() could not be retrieved from DB. failed QueryRow.Scan err: %v", err)
		return 0, err
	}
	return existingMaxId, nil
}

func (P PGX) Exist(id int32) bool {
	P.log.Printf("trace : entering Exist(%d)", id)
	var count int32 = 0

	err := P.con.QueryRow(context.Background(), treesExist, id).Scan(&count)
	if err != nil {
		P.log.Printf("error : Exist(%d) could not be retrieved from DB. failed QueryRow.Scan err: %v", id, err)
		return false
	}
	if count > 0 {
		P.log.Printf("info: Exist(%d) id does exist count:%v", id, count)
		return true
	} else {
		P.log.Printf("info : Exist(%d) id does not exist count:%v", id, count)
		return false
	}
}

func (P PGX) Count() (int32, error) {
	P.log.Println("trace : entering Count()")
	var count int32
	err := P.con.QueryRow(context.Background(), treesCount).Scan(&count)
	if err != nil {
		P.log.Printf("error : Count() could not be retrieved from DB. failed Query.Scan err: %v", err)
		return 0, err
	}
	return count, nil
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
	P.log.Printf("trace : entering Update(%q,%q,%#v)", object.Name, object.Geom, object.TreeAttributes)

	now := time.Now()
	object.LastModificationTime = &now
	if !object.IsActive {
		object.InactivationTime = &now
	} else {
		object.InactivationTime = nil
	}
	P.log.Printf("info : just before Update(%+v)", object)

	res, err := P.con.Exec(context.Background(), treesUpdate, 
		object.Name, &object.Description, &object.ExternalId, object.IsActive, &object.InactivationTime, &object.InactivationReason,
		&object.Comment, &object.IsValidated, &object.IdValidator, &object.LastModificationUser, object.Geom, &object.TreeAttributes, id)	
	if err != nil {
		return nil, GetErrorF("error : Update() query failed", err)
	}
	if res.RowsAffected() < 1 {
		return nil, GetErrorF("error : Update() no row modified", err)
	}
	updatedTree, err := P.Get(id)
	if err != nil {
		return nil, GetErrorF("error : Update() user updated, but cannot be retrieved", err)
	}
	return updatedTree, nil
}

func (P PGX) Delete(id int32) error {
	P.log.Printf("trace : entering Delete(%d)", id)

	res, err := P.con.Exec(context.Background(), treesDelete, id)
	if err != nil {
		return GetErrorF("error : tree could not be deleted", err)
	}
	if res.RowsAffected() < 1 {
		return GetErrorF("error : tree was not deleted", err)
	}

	return nil
}

func (P PGX) SearchTreesByName(pattern string) ([]*TreeList, error) {
	//TODO implement me
	panic("implement me")
}

func (P PGX) IsTreeActive(id int32) bool {
	var isActive bool
	err := P.con.QueryRow(context.Background(), "SELECT is_active FROM tree_mobile WHERE id = $1", id).Scan(&isActive)
	if err != nil {
		P.log.Printf("error : IsTreeActive(%d) could not be retrieved from DB. failed QueryRow.Scan err: %v", id, err)
		return false
	}
	return isActive
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