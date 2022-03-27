package bot

import (
	"context"
	"fmt"
	"regexp"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/firestore"
	"github.com/bwmarrin/discordgo"
	bq "github.com/mager/premintbot/bigquery"
	"go.uber.org/zap"
)

func premintSetRoleCommand(
	ctx context.Context,
	logger *zap.SugaredLogger,
	database *firestore.Client,
	bqClient *bigquery.Client,
	s *discordgo.Session,
	m *discordgo.MessageCreate,
) {
	if m.Content == "!premint-set-role" {
		s.ChannelMessageSend(m.ChannelID, "Missing role. Please use `!premint-set-role <Discord role ID>` to set it. You can find your Discord role ID by going to Server Settings > Roles > Right click the role > Copy ID.")
		return
	}

	// Regex match !premint-set-role <Role ID>
	re := regexp.MustCompile(`^!premint-set-role (.*)$`)
	match := re.FindStringSubmatch(m.Content)

	if len(match) != 2 {
		return
	}

	p := GetConfig(ctx, logger, database, m.GuildID)
	g := getGuildFromMessage(s, m)

	if !isAdmin(p, m.Author) {
		s.ChannelMessageSend(m.ChannelID, "❌ You do not have the Premintbot role. Please contact a server administrator to add it to your account.")
		bq.RecordAdminAction(bqClient, m, "set-role", "not-admin")
		return
	}

	roleID := match[1]
	roleName := ""

	for _, role := range g.Roles {
		if role.ID == roleID {
			roleName = role.Name
		}
	}

	// Make sure the user has the Premintbot role: loop through their roles and make sure they have the guild admin role.
	for _, admin := range p.Config.GuildAdmins {
		if admin == m.Author.ID {
			p.doc.Ref.Update(ctx, []firestore.Update{
				{Path: "premint-role-id", Value: roleID},
				{Path: "premint-role-name", Value: roleName},
			})
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("✅ Premint role updated: %s", roleName))
			bq.RecordAdminAction(bqClient, m, "set-role", "success")
			return
		}
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("❌ Role %s not found. You can find your Discord role ID by going to Server Settings > Roles > Right click the role > Copy ID.", roleID))
	bq.RecordAdminAction(bqClient, m, "set-role", "role-not-found")
}
