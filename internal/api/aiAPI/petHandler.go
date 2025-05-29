package aiAPI

import (
	"encoding/json"
	"net/http"
	"strconv"
	aiService "goAccounting/internal/service/thirdparty/ai"
)

type PetHandler struct {
	petService *aiService.PetService
}

func NewPetHandler() *PetHandler {
	return &PetHandler{
		petService: aiService.NewPetService(),
	}
}

func (h *PetHandler) UpdatePetMood(w http.ResponseWriter, r *http.Request) {
	var goal aiService.FinancialGoal
	if err := json.NewDecoder(r.Body).Decode(&goal); err != nil {
		http.Error(w, "无效的请求数据", http.StatusBadRequest)
		return
	}

	mood, err := h.petService.UpdatePetMood(goal, r.Context())
	if err != nil {
		http.Error(w, "更新宠物心情失败", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mood)
}

func (h *PetHandler) GetDailyEncouragement(w http.ResponseWriter, r *http.Request) {
	encouragement, err := h.petService.GetDailyEncouragement(r.Context())
	if err != nil {
		http.Error(w, "获取每日鼓励失败", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"encouragement": encouragement,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *PetHandler) GetUserPet(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.URL.Query().Get("user_id")
	userId, err := strconv.ParseUint(userIdStr, 10, 32)
	if err != nil {
		http.Error(w, "无效的用户ID", http.StatusBadRequest)
		return
	}

	pet, err := h.petService.GetUserPet(uint(userId))
	if err != nil {
		http.Error(w, "获取宠物信息失败", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pet)
}
