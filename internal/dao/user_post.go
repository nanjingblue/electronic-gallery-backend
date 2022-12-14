package dao

import (
	"electronic-gallery/global"
	"electronic-gallery/internal/model"
)

type userPostDAO struct{}

var UserPostDAO *userPostDAO

func init() {
	UserPostDAO = &userPostDAO{}
}

func (up userPostDAO) GetPostsLikedByUser(uid uint) ([]model.Post, error) {
	var posts []model.Post
	err := global.DBEngine.Table("posts").Select("*").Joins("join user_posts on users.id = user_posts.user_id and liked = ?", true).Scan(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (up userPostDAO) IsLikedByUser(userID, postID uint) bool {
	var userPost model.UserPost
	return global.DBEngine.Where("user_id = ? AND post_id = ? AND liked = ?", userID, postID, 1).First(&userPost).Error == nil
}

func (up userPostDAO) IsCollectedByUser(userID, postID uint) bool {
	var userPost model.UserPost
	return global.DBEngine.Where("user_id = ? AND post_id = ? AND collected = ?", userID, postID, true).First(&userPost).Error == nil
}

func (up userPostDAO) GetUserPost(userID, postID uint) (model.UserPost, error) {
	var userPost model.UserPost
	err := global.DBEngine.Where("user_id = ? AND post_id = ?", userID, postID).First(&userPost).Error
	if err != nil {
		return model.UserPost{}, err
	}
	return userPost, nil
}

func (up userPostDAO) Like(userID, postID uint) error {
	// 先查看是否有记录
	userPost, err := up.GetUserPost(userID, postID)
	if err != nil {
		// 没有记录则插入一个记录
		userPost.UserID = userID
		userPost.PostID = postID
		userPost.Liked = true
		return global.DBEngine.Create(&userPost).Error
	}
	userPost.Liked = true
	return global.DBEngine.Save(&userPost).Error
}

func (up userPostDAO) CancelLike(userID, postID uint) error {
	// 先查看是否有记录
	userPost, err := up.GetUserPost(userID, postID)
	if err != nil {
		// 取消时 一定会先有记录 如果没有直接返回错误
		return err
	}
	userPost.Liked = false
	return global.DBEngine.Save(&userPost).Error
}

func (up userPostDAO) Collect(userID, postID uint) error {
	// 先查看是否有记录
	userPost, err := up.GetUserPost(userID, postID)
	if err != nil {
		userPost.UserID = userID
		userPost.PostID = postID
		userPost.Collected = true
		return global.DBEngine.Create(&userPost).Error
	}
	userPost.Collected = true
	return global.DBEngine.Save(&userPost).Error
}

func (up userPostDAO) CancelCollect(userID, postID uint) error {
	// 先查看是否有记录
	userPost, err := up.GetUserPost(userID, postID)
	if err != nil {
		// 取消时 一定会先有记录 如果没有直接返回错误
		return err
	}
	userPost.Collected = false
	return global.DBEngine.Save(&userPost).Error
}

func (up userPostDAO) Comment(userID, postID uint) error {
	// 先查看是否有记录
	userPost, err := up.GetUserPost(userID, postID)
	if err != nil {
		userPost.UserID = userID
		userPost.PostID = postID
		userPost.Commented = true
		return global.DBEngine.Create(&userPost).Error
	}
	userPost.Commented = true
	return global.DBEngine.Save(&userPost).Error
}

func (up userPostDAO) CancelComment(userID, postID uint) error {
	// 先查看是否有记录
	userPost, err := up.GetUserPost(userID, postID)
	if err != nil {
		// 取消时 一定会先有记录 如果没有直接返回错误
		return err
	}
	userPost.Commented = false
	return global.DBEngine.Save(&userPost).Error
}
