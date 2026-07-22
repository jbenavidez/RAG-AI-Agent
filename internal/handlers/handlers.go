package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"rag/internal/render"
	"rag/internal/services"
)

type Handlers struct {
	uploadService *services.UploadService
	ragService    *services.RagService
	Renderer      *render.Renderer
	wsChan        chan WsMessage
	clients       map[*WebSocketConnection][]string
	mu            sync.Mutex
	MemoryService *services.MemoryService
}

func New(uploadService *services.UploadService, ragService *services.RagService, renderer *render.Renderer, memoryServices *services.MemoryService) *Handlers {
	h := &Handlers{
		uploadService: uploadService,
		ragService:    ragService,
		Renderer:      renderer,
		wsChan:        make(chan WsMessage),
		clients:       make(map[*WebSocketConnection][]string),
		MemoryService: memoryServices,
	}
	go h.ListenToWsChannel() // set go rutine to listen ws chan
	return h
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	if err := h.Renderer.Render(w, "home.html", nil); err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) WsChat(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println("unable to start ws", err)
		return
	}
	//set ws res
	var response WsJsonResponse
	response.Message = `<em><small> connected to served</small></em>`
	conn := &WebSocketConnection{Conn: ws}

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err)
	}
	//
	log.Println("cconnected success ")

	go h.ListenForWs(conn) // start go runtine to listen Ws

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

		h.wsChan <- wsMessage
	}
}

func (h *Handlers) ListenToWsChannel() {
	for e := range h.wsChan {
		fmt.Println("listening for ws event")

		switch e.Payload.Action {
		case "ask":
			answer, err := h.ragService.AskQuestion(e.Payload.Message)
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

			fmt.Println("user question", e.Payload.Message)
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

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "missing uploaded file", http.StatusBadRequest)
		return
	}
	defer func() {
		_ = file.Close()
	}()

	description := r.FormValue("description")

	_, err = h.uploadService.SaveUploadedFile(r.Context(), file, header, description)
	if err != nil {
		http.Error(w, "unable to save uploaded file", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/docs", http.StatusSeeOther)
}
