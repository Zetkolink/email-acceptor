package models

import "github.com/jinzhu/gorm"

type Messages struct {
	UniqueId string
	Obj      string
}

func (m Messages) Migrate(db *gorm.DB) {
	db.Exec("create or replace view v_messages as select unique_id as id,json_build_object('id',unique_id,'sender',sender,'created_at',created_at,'to',json_agg(messages.to),'subject',subject,'message',message,'state',state) as obj from messages group by unique_id, sender, created_at, subject, message, state;")
}

func (m Messages) TableName() string {
	return "v_messages"
}
