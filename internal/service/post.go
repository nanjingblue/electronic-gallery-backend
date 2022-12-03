package service

import (
	"electronic-gallery/internal/dao"
	"electronic-gallery/internal/model"
	"electronic-gallery/internal/serializer"
)

type PostService struct{}

type PostGetListService struct {
	PostService
}

func (p PostGetListService) GetList(svc *Service) serializer.Response {
	u, _ := svc.ctx.Get("user")
	user := u.(model.User)

	posts, err := dao.Post.GetPosts(user.ID)
	if err != nil {
		return serializer.Response{
			Code:  500,
			Msg:   "服务器内部错误",
			Error: err.Error(),
		}
	}

	return serializer.Response{
		Code: 200,
		Data: serializer.BuildPostsWithMe(posts, user),
		Msg:  "获取所有follow post 成功",
	}
}

type PostGetMyListService struct{}

func (p *PostGetMyListService) GetList(svc *Service) serializer.Response {
	u, _ := svc.ctx.Get("user")
	user := u.(model.User)

	posts, err := dao.Post.GetAllPostByUserID(user.ID)
	if err != nil {
		return serializer.Response{
			Code:  500,
			Msg:   "服务器内部错误",
			Error: err.Error(),
		}
	}

	return serializer.Response{
		Code: 200,
		Data: serializer.BuildPosts(posts),
		Msg:  "获取所有my post 成功",
	}
}

type PostCreateService struct {
	PostService
	Content string `form:"content" json:"content" binding:"required"`
	Path    string `form:"path" json:"path"`
}

func (p *PostCreateService) CreatePost(svc *Service) serializer.Response {
	u, _ := svc.ctx.Get("user")
	user := u.(model.User)

	post := model.Post{
		Content: p.Content,
		UserID:  user.ID,
		Image:   p.Path,
	}

	err := dao.Post.CreatePost(&post)
	if err != nil {
		return serializer.Response{
			Code:  500,
			Msg:   "服务器内部错误",
			Error: err.Error(),
		}
	}

	return serializer.Response{
		Code: 200,
		Data: serializer.BuildPost(&post),
		Msg:  "创建post成功",
	}
}

type PostLikeService struct {
	PostID uint `form:"post_id" json:"post_id" binding:"required"`
}

func (p *PostLikeService) Like(svc *Service) serializer.Response {
	u, _ := svc.ctx.Get("user")
	user := u.(model.User)

	post, err := dao.Post.GetPostByID(p.PostID)
	if err != nil {
		return serializer.Response{
			Code:  400,
			Msg:   "like post: 不存在post",
			Error: err.Error(),
		}
	}

	post.AddLike()                               // redis
	err = dao.UserPostDAO.Like(user.ID, post.ID) // 关系
	if err != nil {
		return serializer.Response{
			Code:  500,
			Msg:   "like post: 添加记录失败",
			Error: err.Error(),
		}
	}

	return serializer.Response{
		Code: 200,
		Data: serializer.BuildPost(&post),
		Msg:  "like post success",
	}
}

func (p *PostLikeService) CancelLike(svc *Service) serializer.Response {
	u, _ := svc.ctx.Get("user")
	user := u.(model.User)

	post, err := dao.Post.GetPostByID(p.PostID)
	if err != nil {
		return serializer.Response{
			Code:  400,
			Msg:   "cancel like post: 不存在post",
			Error: err.Error(),
		}
	}
	post.CancelLike()
	err = dao.UserPostDAO.CancelLike(user.ID, post.ID)
	if err != nil {
		return serializer.Response{
			Code:  500,
			Msg:   "cancel like post: 不存在post 或者 更新失败",
			Error: err.Error(),
		}
	}

	return serializer.Response{
		Code: 200,
		Data: post,
		Msg:  "cancel like post success",
	}
}

type PostCollectionService struct {
	PostID uint `form:"post_id" json:"post_id" binding:"required"`
}

func (p *PostCollectionService) Collection(svc *Service) serializer.Response {
	u, _ := svc.ctx.Get("user")
	user := u.(model.User)

	post, err := dao.Post.GetPostByID(p.PostID)
	if err != nil {
		return serializer.Response{
			Code:  400,
			Msg:   "collect post: 不存在post",
			Error: err.Error(),
		}
	}

	post.AddCollection()
	err = dao.UserPostDAO.Collect(user.ID, post.ID)
	if err != nil {
		return serializer.Response{
			Code:  500,
			Msg:   "collect post: 添加后更新记录失败",
			Error: err.Error(),
		}
	}

	return serializer.Response{
		Code: 200,
		Data: serializer.BuildPost(&post),
		Msg:  "collect post success",
	}
}

func (p *PostCollectionService) CancelCollection(svc *Service) serializer.Response {
	u, _ := svc.ctx.Get("user")
	user := u.(model.User)

	post, err := dao.Post.GetPostByID(p.PostID)
	if err != nil {
		return serializer.Response{
			Code:  400,
			Msg:   "cancel collect post: 不存在post",
			Error: err.Error(),
		}
	}

	post.CancelCollection()
	err = dao.UserPostDAO.CancelCollect(user.ID, post.ID)
	if err != nil {
		return serializer.Response{
			Code:  500,
			Msg:   "cancel collect failed: 不存在post或更新失败",
			Error: err.Error(),
		}
	}

	return serializer.Response{
		Code: 200,
		Data: post,
		Msg:  "cancel like post success",
	}
}
