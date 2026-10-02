package models

import "time"

type Komentar struct {
	ID        int     `json:"id"`
	BeritaID  int     `json:"beritaId"`
	UserID    int     `json:"userId"`
	Isi       string  `json:"isi"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	//Data user yang membuat komentar
	UserNama string `json:"userNama,omitempty"`
}