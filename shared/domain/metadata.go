package domain

import "github.com/vaxxnsh/shared/db"

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

type UserAvatar struct {
	UserID   string
	ImageUrl string
	Name     string
}

func MapBulkUserAvatarsToDomain(rows []db.GetBulkUserAvatarsRow) []UserAvatar {
	result := make([]UserAvatar, 0, len(rows))
	for _, r := range rows {
		result = append(result, UserAvatar{
			UserID:   r.UserID.String(),
			ImageUrl: r.ImageUrl,
			Name:     r.AvatarName.String,
		})
	}
	return result
}
