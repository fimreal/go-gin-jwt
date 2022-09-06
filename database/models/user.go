/*
user 表

CREATE TABLE `user` (

	`user_id` bigint NOT NULL AUTO_INCREMENT,
	`username` varchar(32) NOT NULL,
	`password` varchar(256) NOT NULL,
	`email` varchar(128) DEFAULT NULL,
	`level` int(3) DEFAULT 0 NOT NULL,
	PRIMARY KEY (`user_id`),
	UNIQUE KEY `email` (`email`)
	) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4;
*/
package models

import "database/sql"

type User struct {
	UserID   int            `gorm:"primaryKey;column:user_id;" json:"user_id"`
	Username string         `gorm:"column:username" json:"username"`
	Password string         `gorm:"column:password" json:"password"`
	Email    sql.NullString `gorm:"column:email" json:"email"`
	Level    int            `gorm:"column:level" json:"level"`
}

type AddUserData struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required"`
}

type LoginRequest struct {
	// UserID   int    `json:"user_id"`
	// Username string `json:"username"`
	User     string `json:"user"`
	Password string `json:"password" validate:"required"`
}

type LoginResult struct {
	UserID       int    `json:"user_id"`
	Username     string `json:"username"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
