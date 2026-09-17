package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	qrcode "github.com/skip2/go-qrcode"
)

// GetTicketQR renders a PNG QR code that encodes the ticket/booking code in the
// URL. The frontend embeds it with a plain <img> tag, so no client-side QR
// library is needed. The QR only carries the code — a pointer to the ticket in
// the database — never the ticket data itself.
func (h *Handler) GetTicketQR(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		http.Error(w, "code required", http.StatusBadRequest)
		return
	}

	png, err := qrcode.Encode(code, qrcode.Medium, 256)
	if err != nil {
		http.Error(w, "could not generate qr", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	_, _ = w.Write(png)
}
