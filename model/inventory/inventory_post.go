package inventory

import (
	"fmt"

	"github.com/micro/go-micro/util/log"
	"github.com/xiaobudongzhang/micro-basic/common"
	"github.com/xiaobudongzhang/micro-basic/db"
	proto "github.com/xiaobudongzhang/micro-inventory-srv/proto/inventory"
)

func (s *service) Sell(bookId int64, userId int64) (id int64, err error) {
	tx, err := db.GetDB().Begin()
	if err != nil {
		log.Logf("事务开启失败", err.Error())
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	querySQL := `select id,book_id,unit_price,stock,version from inventory where book_id = ?`

	inv := &proto.Inv{}

	updateSQL := `update inventory set stock = ?, version = ? where book_id = ? and version = ?`

	var deductInv func() error

	deductInv = func() (errIn error) {
		errIn = tx.QueryRow(querySQL, bookId).Scan(&inv.Id, &inv.BookId, &inv.UnitPrice, &inv.Stock, &inv.Version)
		if errIn != nil {
			log.Logf("查询失败 %s", errIn)
			return errIn
		}

		if inv.Stock < 1 {
			errIn = fmt.Errorf("库存不足")
			log.Logf(errIn.Error())
			return errIn
		}

		r, errIn := tx.Exec(updateSQL, inv.Stock-1, inv.Version+1, bookId, inv.Version)
		if errIn != nil {
			log.Logf("更新数据库失败 %s", errIn)
			return
		}

		if affected, _ := r.RowsAffected(); affected == 0 {
			log.Logf("更新失败 %d", inv.Version)
			deductInv()
		}
		return
	}
	// 开始销存
	err = deductInv()
	if err != nil {
		log.Logf("[Sell] 销存失败，err：%s", err)
		return
	}

	insertSQL := `insert inventory_history (book_id,user_id,state) value (?, ?, ?)`
	r, err := tx.Exec(insertSQL, bookId, userId, common.InventoryHistoryStateNotOut)
	if err != nil {
		log.Logf("新增销存记录失败 %s", err)
		return
	}
	id, _ = r.LastInsertId()

	tx.Commit()
	return
}

func (s *service) Confirm(id int64, state int) (err error) {
	updateSQL := `update inventory_history set state = ? where id = ?;`

	o := db.GetDB()

	_, err = o.Exec(updateSQL, state, id)
	if err != nil {
		log.Logf("confirm 更新失败 %s", err)
		return
	}
	return
}
