package model

type Message struct {
	BaseModel
	FromUserID uint   `json:"from_user_id" gorm:"index"`
	ToUserID   uint   `json:"to_user_id" gorm:"index"`      // 私聊对象，群聊为0
	RoomID     string `json:"room_id" gorm:"index;size:64"` // 房间ID，私聊为空
	Content    string `json:"content" gorm:"type:text"`
	MsgType    string `json:"msg_type" gorm:"size:20;default:text"` // text / image /system
	FromUser   *User  `json:"from_user,omitempty" gorm:"foreignKey:FromUserID;references:ID"`
}
