package inventory

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/xiaobudongzhang/seata-golang/client/at/exec"
	"github.com/xiaobudongzhang/seata-golang/client/config"
	"github.com/xiaobudongzhang/seata-golang/client/context"

	"github.com/micro/go-micro/util/log"
	"github.com/xiaobudongzhang/micro-basic/common"
	proto "github.com/xiaobudongzhang/micro-inventory-srv/proto/inventory"
	"github.com/xiaobudongzhang/micro-plugins/db"
)
func NextSnowflakeId() uint64 {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
	return uint64(node.Generate())
}
func (s *service) Sell(bookId int64, userId int64, ctx2 *context.RootContext) (id int64, err error) {
	log.Log("Sell")
	//tx, err := db.GetDB().Begin()



	db,err := exec.NewDB(config.GetClientConfig().ATConfig)
	if err != nil {
		panic(err)
	}
	tx2, err2 := db.Begin(ctx2)
	if err2 != nil {
		log.Logf("事务开启失败", err.Error())
		return
	}
	defer func() {
		if err != nil {
			tx2.Rollback()
		}
	}()

	querySQL := `select id,book_id,unit_price,stock,version from micro_book_mall.inventory where book_id = ?`

	inv := &proto.Inv{}

	updateSQL := `update micro_book_mall.inventory set stock = ?, version = ? where book_id = ? and version = ?`

	var deductInv func() error

	deductInv = func() (errIn error) {
		errIn = db.QueryRow(querySQL, bookId).Scan(&inv.Id, &inv.BookId, &inv.UnitPrice, &inv.Stock, &inv.Version)
		if errIn != nil {
			log.Logf("查询失败 %s", errIn)
			return errIn
		}

		if inv.Stock < 1 {
			errIn = fmt.Errorf("库存不足")
			log.Logf(errIn.Error())
			return errIn
		}

		r, errIn := db.Exec(updateSQL, inv.Stock-1, inv.Version+1, bookId, inv.Version)
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

	insertSQL := `insert into micro_book_mall.inventory_history (id,book_id,user_id,state) value (?,?, ?, ?)`
	r, err := db.Exec(insertSQL,NextSnowflakeId(), bookId, userId, common.InventoryHistoryStateNotOut)
	if err != nil {
		log.Logf("新增销存记录失败 %s", err)
		return
	}
	id, _ = r.LastInsertId()

	tx2.Commit()
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
