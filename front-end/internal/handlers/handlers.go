package handlers

import (
	"fmt"
	"frontend/internal/config"
	"frontend/internal/models"
	"frontend/internal/render"
	pb "frontend/proto/generated"
	"net/http"
)

type Repository struct {
	App *config.AppConfig
}

var Repo *Repository

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

// ChatRoom is the handler for the chatroom page
func (m *Repository) ChatRoom(w http.ResponseWriter, r *http.Request) {

	// Render template
	render.RenderTemplate(w, r, "chatroom.page.tmpl", &models.TemplateData{
		Data: nil,
	})
}

// AnswerUserQuestion: answer user question.
func (m *Repository) AnswerUserQuestion(w http.ResponseWriter, r *http.Request) {
	var payload models.Userquestion

	err := m.readJSON(w, r, &payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	//set grpc req
	req := &pb.AIAgentRequest{
		Question: payload.Question,
	}

	// call grcp
	response, err := m.App.GRPCClient.GetAIAgentAnswerFromUserQuestion(r.Context(), req)

	if err != nil {
		// set fail resposne
		errRes := models.JSONResponse{
			Error:   true,
			Message: "Server Error, Unable to answer the question",
		}
		_ = m.writeJSON(w, http.StatusOK, errRes)
	}

	fmt.Println("Valinor_calling....", response.Answer)
	// set response
	agentAnwer := models.AgentResponse{
		Answer: response.Answer,
	}
	//set response
	resp := models.JSONResponse{
		Error:   false,
		Message: "success",
		Data:    agentAnwer,
	}
	// send resposne
	fmt.Println("the response", resp)
	_ = m.writeJSON(w, http.StatusOK, resp)

}
