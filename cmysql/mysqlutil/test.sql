#order
INSERT INTO `test`.`order`(`oid`, `name`, `ct`) VALUES(?, ?, ?) ON DUPLICATE KEY UPDATE `name`=?, `ct`=?;
SELECT `oid`, `name`, `ct` FROM `test`.`order` WHERE `oid`=?;
UPDATE `test`.`order` SET `oid`=?, `name`=?, `ct`=? WHERE `oid`=?;
type Order struct {
	Oid int `json:"oid,omitempty"`
	Name string `json:"name,omitempty"`
	Ct string `json:"ct,omitempty"`
}

#user
INSERT INTO `test`.`user`(`id`, `name`, `ct`) VALUES(?, ?, ?) ON DUPLICATE KEY UPDATE `ct`=?;
SELECT `id`, `name`, `ct` FROM `test`.`user` WHERE `id`=? AND `name`=?;
UPDATE `test`.`user` SET `id`=?, `name`=?, `ct`=? WHERE `id`=? AND `name`=?;
type User struct {
	Id int `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Ct string `json:"ct,omitempty"`
}

