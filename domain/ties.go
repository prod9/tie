package domain

import (
	"context"
	"regexp"
	"time"

	"tie.prodigy9.co/data"
)

const SlugRx = `[0-9a-zA-Z\-]+`

func ValidSlug(v string) error {
	if ok, _ := regexp.MatchString("^"+SlugRx+"$", v); ok {
		return nil
	}

	return &ErrValidation{
		Field:   "slug",
		Message: "must contain only A-Z and 0-9",
	}
}

type Tie struct {
	ID        int    `json:"id" db:"id"`
	Slug      string `json:"slug" db:"slug"`
	TargetURL string `json:"target_url" db:"target_url"`

	CreatedAt time.Time `json:"created_at" db:"ctime"`
	UpdatedAt time.Time `json:"updated_at" db:"mtime"`
}

func (t *Tie) String() string {
	return t.Slug + " => " + t.TargetURL
}

func ListAllTies(ctx context.Context, out *List[*Tie]) error {
	if out.Data == nil {
		out.Data = []*Tie{}
	}

	return data.Select(ctx, &out.Data,
		`SELECT * FROM ties ORDER BY id ASC`)
}

func GetTieBySlug(ctx context.Context, out *Tie, slug string) error {
	if err := Required("slug", slug); err != nil {
		return err
	} else if err = ValidSlug(slug); err != nil {
		return err
	}

	return data.Get(ctx, out,
		`SELECT * FROM ties WHERE slug = $1 ORDER BY id ASC`,
		slug)
}

type CreateTie struct {
	Slug      string `json:"slug"`
	TargetURL string `json:"target_url"`
}

func (t *CreateTie) Validate() error {
	if err := Required("slug", t.Slug); err != nil {
		return err
	}
	if err := ValidSlug(t.Slug); err != nil {
		return err
	}
	if err := Required("target_url", t.TargetURL); err != nil {
		return err
	}
	return nil
}

func (t *CreateTie) Execute(ctx context.Context, out *Tie) error {
	return data.Get(ctx, out, `
    INSERT INTO ties (slug, target_url)
    VALUES ($1, $2)
    ON CONFLICT (slug) DO UPDATE
    SET target_url = $2
    RETURNING *`,
		t.Slug,
		t.TargetURL)
}

type DeleteTie struct {
	Slug string `json:"slug"`
}

func (d *DeleteTie) Validate() error {
	if err := Required("slug", d.Slug); err != nil {
		return err
	} else if err = ValidSlug(d.Slug); err != nil {
		return err
	}
	return nil
}

func (d *DeleteTie) Execute(ctx context.Context, out *Tie) error {
	return data.Get(ctx, out, `DELETE FROM ties WHERE slug = $1 RETURNING *`,
		d.Slug)
}
