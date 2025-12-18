package handlers

import (
	"context"
	"fmt"
	"frontend/internal/config"
	pb "frontend/proto/generated"
	"log"
)

var app *config.AppConfig

func NewHelpers(a *config.AppConfig) {
	app = a
}

type WsPayload struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

type WsMessage struct {
	Payload *WsPayload           // the JSON payload
	Conn    *WebSocketConnection // pointer to the live connection
}

// ListenForWs : listen for user question
func ListenForWs(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Errror", fmt.Sprintf("%v", r))
		}
	}()

	var payload WsPayload

	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			// do nothing
			log.Println("ws read err", err)
			break
		} else {
			wsMessage := WsMessage{
				Payload: &payload,
				Conn:    conn,
			}
			fmt.Println("Sending payload to channel", wsMessage)
			wsChan <- wsMessage // send payload to channel
		}
	}
}

// ListenToWsChannel : listen to channel --> get question and broadcast and answer
func ListenToWsChannel() {
	var response WsJsonResponse
	for {
		e := <-wsChan // read payload from channel
		fmt.Println("listning fo webhook event")
		switch e.Payload.Action {

		case "question":

			//TODO: get response from the agent

			// Get answer for user question from Agent
			answer, err := SendQuestionToAgent(e.Payload.Message)
			if err != nil {
				fmt.Println("unable to get answer from agent", err)
				break
			}
			// store answer
			clients[e.Conn] = append(clients[e.Conn], answer)
			//send response back
			response.Action = "answer"
			response.Message = answer

			fmt.Println("gondor about to response", response)
			BroadcastResponseToConn(e.Conn, response)

		}
	}
}

func SendQuestionToAgent(question string) (string, error) {
	// set grpc req
	req := &pb.AIAgentRequest{
		Question: question,
	}
	fmt.Println("valinor_question", question)
	// call grcp
	response, err := app.GRPCClient.GetAIAgentAnswerFromUserQuestion(context.Background(), req)
	if err != nil {
		return "", err
	}
	return response.Answer, nil
}

func BroadcastResponseToConn(conn *WebSocketConnection, response WsJsonResponse) {

	if conn == nil || conn.Conn == nil {
		log.Println("nil websocket connection")
		return
	}

	err := conn.WriteJSON(response)
	if err != nil {
		log.Println("WS err", err)
		_ = conn.Close()
	}
}
