// ai理财宠物会根据当前用户是否完成理财目标转化心情，
// 需要和前端交互的是一个标识心情的词语和一句简短的鼓励用户的话，
// 通过调用云端大模型api实现
package aiService

import (
	"context"
	"math/rand"
	"time"
)

type PetService struct{}

type PetMood struct {
	Mood          string    `json:"mood"`
	Encouragement string    `json:"encouragement"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type FinancialGoal struct {
	Target    float64 `json:"target"`
	Current   float64 `json:"current"`
	Completed bool    `json:"completed"`
}

func NewPetService() *PetService {
	return &PetService{}
}

func (ps *PetService) UpdatePetMood(goal FinancialGoal, ctx context.Context) (*PetMood, error) {
	// 根据理财目标完成情况调用AI生成心情和鼓励语
	completionRate := goal.Current / goal.Target

	var mood string
	var encouragements []string

	switch {
	case completionRate >= 1.0:
		mood = "开心"
		encouragements = []string{
			"太棒了！你已经完成了理财目标！",
			"恭喜你达成目标，继续保持这种好习惯！",
			"你的坚持得到了回报，为你感到骄傲！",
		}
	case completionRate >= 0.8:
		mood = "满意"
		encouragements = []string{
			"你做得很好，离目标越来越近了！",
			"继续努力，成功就在眼前！",
			"你的进步让我很欣慰，加油！",
		}
	case completionRate >= 0.5:
		mood = "普通"
		encouragements = []string{
			"还不错哦，继续朝着目标前进吧！",
			"每一小步都是进步，不要放弃！",
			"相信自己，你一定可以达成目标的！",
		}
	default:
		mood = "担心"
		encouragements = []string{
			"别灰心，制定一个小目标开始吧！",
			"理财需要耐心，慢慢来不要急！",
			"每天进步一点点，我会陪着你的！",
		}
	}

	// 随机选择一句鼓励语
	encouragement := encouragements[rand.Intn(len(encouragements))]

	return &PetMood{
		Mood:          mood,
		Encouragement: encouragement,
		UpdatedAt:     time.Now(),
	}, nil
}

func (ps *PetService) GetDailyEncouragement(ctx context.Context) (string, error) {
	// 获取每日鼓励语
	encouragements := []string{
		"新的一天，新的开始！今天也要好好管理财务哦！",
		"记得记录今天的收支，小宠物在等着看你的进步呢！",
		"理财路上有我陪伴，一起加油吧！",
		"今天也要为了梦想而努力存钱哦！",
	}

	encouragement := encouragements[rand.Intn(len(encouragements))]
	return encouragement, nil
}

func (ps *PetService) GetUserPet(userId uint) (*PetMood, error) {
	// TODO: Implement database query
	return &PetMood{
		Mood:          "开心",
		Encouragement: "你今天做得很棒！",
		UpdatedAt:     time.Now(),
	}, nil
}
