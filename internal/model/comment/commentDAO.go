package commentModel

import (
	"errors"
	"goAccounting/global"
	transactionModel "goAccounting/internal/model/transaction"
	userModel "goAccounting/internal/model/user"

	"gorm.io/gorm"
)

type CommentDao struct {
	db *gorm.DB
}

func NewDao(db ...*gorm.DB) *CommentDao {
	if len(db) > 0 {
		return &CommentDao{db: db[0]}
	}
	return &CommentDao{global.GlobalDb}
}

type CommentCreateData struct {
	UserId        uint
	TargetUserId  uint
	TransactionId uint
	Content       string
}

func (cd *CommentDao) Create(data CommentCreateData) (Comment, error) {
	if err := cd.CheckCommentPermission(data.UserId, data.TargetUserId, data.TransactionId); err != nil {
		return Comment{}, err
	}

	comment := Comment{
		UserId:        data.UserId,
		TargetUserId:  data.TargetUserId,
		TransactionId: data.TransactionId,
		Content:       data.Content,
	}

	err := cd.db.Create(&comment).Error
	return comment, err
}

func (cd *CommentDao) CheckCommentPermission(userId, targetUserId, transactionId uint) error {
	userDao := userModel.NewDao(cd.db)
	isFriend, err := userDao.IsRealFriend(userId, targetUserId)
	if err != nil {
		return err
	}
	if !isFriend {
		return errors.New("not friend relationship")
	}

	var config userModel.TransactionShareConfig
	err = cd.db.Where("user_id = ?", targetUserId).First(&config).Error
	if err != nil {
		return err
	}
	if !config.IsShared {
		return errors.New("target user has not shared transactions")
	}

	var transaction transactionModel.Transaction
	err = cd.db.Where("id = ? AND user_id = ?", transactionId, targetUserId).First(&transaction).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("transaction not found or not owned by target user")
		}
		return err
	}

	return nil
}

func (cd *CommentDao) Delete(commentId, userId uint) error {
	result := cd.db.Where("id = ? AND user_id = ?", commentId, userId).Delete(&Comment{})
	if result.RowsAffected == 0 {
		return errors.New("comment not found or not owned by the user")
	}
	return result.Error
}

type CommentUpdateData struct {
	Content string
}

func (cd *CommentDao) Update(commentId, userId uint, data CommentUpdateData) error {
	result := cd.db.Model(&Comment{}).
		Where("id = ? AND user_id = ?", commentId, userId).
		Update("content", data.Content)

	if result.RowsAffected == 0 {
		return errors.New("comment not found or not owned by the user")
	}
	return result.Error
}

type ListOptions struct {
	UserId        *uint  // 评论者ID
	TargetUserId  *uint  // 目标用户ID
	TransactionId *uint  // 交易记录ID
	OrderBy       string // 排序字段
	OrderDesc     bool   // 是否降序
	Limit         int    // 每页条数
	Offset        int    // 偏移量
}

func (cd *CommentDao) List(options ListOptions) ([]Comment, error) {
	query := cd.db.Model(&Comment{})

	if options.UserId != nil {
		query = query.Where("user_id = ?", *options.UserId)
	}
	if options.TargetUserId != nil {
		query = query.Where("target_user_id = ?", *options.TargetUserId)
	}
	if options.TransactionId != nil {
		query = query.Where("transaction_id = ?", *options.TransactionId)
	}

	orderBy := "created_at"
	if options.OrderBy != "" {
		orderBy = options.OrderBy
	}

	if options.OrderDesc {
		query = query.Order(orderBy + " DESC")
	} else {
		query = query.Order(orderBy)
	}

	if options.Limit > 0 {
		query = query.Limit(options.Limit)
	}

	if options.Offset > 0 {
		query = query.Offset(options.Offset)
	}

	var comments []Comment
	err := query.Find(&comments).Error
	return comments, err
}

func (cd *CommentDao) GetTransactionComments(transactionId uint) ([]Comment, error) {
	var comments []Comment
	err := cd.db.Where("transaction_id = ?", transactionId).
		Order("created_at DESC").
		Find(&comments).Error
	return comments, err
}

func (cd *CommentDao) GetCommentById(commentId uint) (Comment, error) {
	var comment Comment
	err := cd.db.First(&comment, commentId).Error
	return comment, err
}

func (cd *CommentDao) CountByTransaction(transactionId uint) (int64, error) {
	var count int64
	err := cd.db.Model(&Comment{}).
		Where("transaction_id = ?", transactionId).
		Count(&count).Error
	return count, err
}
