package domain

import (
	"context"
	"time"
)

type Entities struct {
	Id        int64     `json:"id" pg:"id,pk"`
	CreatedAt time.Time `json:"created_at" pg:"default:now()"`
	UpdatedAt time.Time `json:"updated_at" pg:"default:now()"`
}

func (e *Entities) BeforeUpdate(ctx context.Context) (context.Context, error) {
	e.UpdatedAt = time.Now()
	return ctx, nil
}
