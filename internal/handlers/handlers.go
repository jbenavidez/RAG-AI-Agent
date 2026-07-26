package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"rag/internal/middleware"
	"rag/internal/render"
	"rag/internal/services"
)

type Handlers struct {
	uploadService   *services.UploadService
	ragService      *services.RagService
	Renderer        *render.Renderer
	chatMessageChan chan WsMessage
	clients         map[*WebSocketConnection][]string
	mu              sync.Mutex
	MemoryService   *services.MemoryService
}

func New(uploadService *services.UploadService, ragService *services.RagService, renderer *render.Renderer, memoryServices *services.MemoryService) *Handlers {
	h := &Handlers{
		uploadService:   uploadService,
		ragService:      ragService,
		Renderer:        renderer,
		chatMessageChan: make(chan WsMessage),
		clients:         make(map[*WebSocketConnection][]string),
		MemoryService:   memoryServices,
	}
	go h.ListenToWsChannel()           // set go rutine to listen ws chan
	go h.ListenToUploadStatusChannel() // set go rutine to  upload status chan
	return h
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	if err := h.Renderer.Render(w, "home.html", nil); err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) WebSocket(w http.ResponseWriter, r *http.Request) {
	// Get SessionID from request context.
	sessionID := middleware.GetSessionID(r)
	if sessionID == "" {
		http.Error(w, "missing session", http.StatusBadRequest)
		return
	}
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println("unable to start ws", err)
		return
	}
	conn := &WebSocketConnection{
		Conn:      ws,
		SessionID: sessionID,
	}
	// Register WebSocket connection immediately.
	h.mu.Lock()
	h.clients[conn] = []string{}
	h.mu.Unlock()

	response := WsJsonResponse{
		Action:    "connected",
		Message:   "Connected to server.",
		SessionID: sessionID,
	}

	if err := ws.WriteJSON(response); err != nil {
		log.Println(err)
		_ = ws.Close()
		return
	}

	log.Println("connected success")

	go h.ListenForWs(conn)
}

func (h *Handlers) ListenForWs(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()

	var payload WsPayload
	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			log.Println("ws read err", err)
			break
		}
		wsMessage := WsMessage{
			Payload: &payload,
			Conn:    conn,
		}

		fmt.Println("Sending  to channel", wsMessage)

		h.chatMessageChan <- wsMessage
	}
}

func (h *Handlers) ListenToWsChannel() {
	for e := range h.chatMessageChan {
		fmt.Println("listening for chatMessageChan event")

		switch e.Payload.Action {
		case "ask":
			chatHistory, err := h.MemoryService.GetChatsHistory(context.Background(), e.Conn.SessionID)
			if err != nil {
				fmt.Println("unable to get chat history", err)
			}

			answer, err := h.ragService.AskQuestion(e.Payload.Message, chatHistory)
			if err != nil {
				fmt.Println("unable to get answer from agent", err)
				response := WsJsonResponse{
					Action:  "error",
					Message: "unable to get answer from agent",
				}
				h.BroadcastResponseToConn(e.Conn, response)
				continue
			}

			h.mu.Lock()
			h.clients[e.Conn] = append(h.clients[e.Conn], answer)
			h.mu.Unlock()
			response := WsJsonResponse{
				Action:  "answer",
				Message: answer,
			}
			// save chat turn
			if err := h.MemoryService.Store(context.Background(), e.Conn.SessionID, e.Payload.Message, answer); err != nil {
				fmt.Println("unable to save chat memory", err)
			}
			h.BroadcastResponseToConn(e.Conn, response)

		default:
			response := WsJsonResponse{
				Action:  "error",
				Message: "unknown websocket action",
			}

			h.BroadcastResponseToConn(e.Conn, response)
		}
	}
}

func (h *Handlers) BroadcastResponseToConn(conn *WebSocketConnection, response WsJsonResponse) {

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

func (h *Handlers) GetAllUploadFiles(w http.ResponseWriter, r *http.Request) {
	files, err := h.uploadService.GetAllUploadedFile()
	if err != nil {
		http.Error(w, "unable to retrieve all uploaded files", http.StatusInternalServerError)
		return
	}

	if err := h.Renderer.Render(w, "all_uploaded_file.html", files); err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) UploadDoc(w http.ResponseWriter, r *http.Request) {
	if err := h.Renderer.Render(w, "upload.html", nil); err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) ProcessDoc(w http.ResponseWriter, r *http.Request) {
	// Limit request size to 10 MB.
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "file too large or invalid form", http.StatusBadRequest)
		return
	}

	// Get SessionID from cokoies
	sessionID := middleware.GetSessionID(r)
	if sessionID == "" {
		http.Error(w, "missing session", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "missing uploaded file", http.StatusBadRequest)
		return
	}
	defer func() {
		_ = file.Close()
	}()

	description := r.FormValue("description")

	_, err = h.uploadService.SaveUploadedFile(r.Context(), file, header, description, sessionID)
	if err != nil {
		http.Error(w, "unable to save uploaded file", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/docs", http.StatusSeeOther)
}

func (h *Handlers) ListenToUploadStatusChannel() {
	for event := range h.uploadService.UploadStatusChan {
		if event == nil {
			continue
		}

		response := WsJsonResponse{
			Action:  "upload_status",
			Message: event.Message,
		}

		h.SendUploadStatusToSession(event.SessionID, response)
	}
}

func (h *Handlers) SendUploadStatusToSession(sessionID string, response WsJsonResponse) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for conn := range h.clients {
		if conn.SessionID == sessionID {
			h.BroadcastResponseToConn(conn, response)
			return
		}
	}
}
