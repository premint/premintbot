package bot

import (
	"context"
	"fmt"
	"regexp"

	"cloud.google.com/go/firestore"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

func premintSetRoleCommand(
	ctx context.Context,
	logger *zap.SugaredLogger,
	database *firestore.Client,
	s *discordgo.Session,
	m *discordgo.MessageCreate,
) {
	if m.Content == "!premint-set-role" {
		s.ChannelMessageSend(m.ChannelID, "Missing role. Please use `!premint-set-role <role name>` to set it. You can find your role ID by right-clicking on the role > Copy ID.")
		return
	}

	// Regex match !premint-set-role <API Key>
	re := regexp.MustCompile(`^!premint-set-role (.*)$`)
	match := re.FindStringSubmatch(m.Content)

	if len(match) != 2 {
		return
	}

	p := GetConfig(ctx, logger, database, m.GuildID)
	g := getGuildFromMessage(s, m)
	roleID := match[1]

	// Make sure the user has the Premintbot role: loop through their roles and make sure they have the guild admin role.
	for _, r := range m.Member.Roles {
		if r == p.Config.GuildAdminRoleID {
			for _, role := range g.Roles {
				if role.ID == roleID {
					p.doc.Ref.Update(ctx, []firestore.Update{
						{Path: "premint-role-id", Value: roleID},
						{Path: "premint-role-name", Value: role.Name},
					})
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("✅ Premint role updated: %s", role.Name))
					return
				}
			}
		}
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("❌ Role %s not found. You can find your role ID but right clicking it > Copy ID.", roleID))
}
