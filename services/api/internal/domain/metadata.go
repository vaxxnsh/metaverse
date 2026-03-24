package domain

import "github.com/vaxxnsh/metaverse/api/internal/db"

type Avatar struct {
	ID       string
	ImageUrl string
	Name     string
}

func DBAvatarToDomain(avatar db.Avatar) *Avatar {
	return &Avatar{
		ID:       avatar.ID.String(),
		Name:     avatar.Name.String,
		ImageUrl: avatar.ImageUrl,
	}
}

func MapAvatarsToDomain(avatars []db.Avatar) []Avatar {
	domainAvatars := make([]Avatar, 0, len(avatars))

	for _, a := range avatars {
		domainAvatar := DBAvatarToDomain(a)
		domainAvatars = append(domainAvatars, *domainAvatar)
	}

	return domainAvatars
}
