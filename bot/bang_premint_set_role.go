package bot

import (
	"context"
	"regexp"

	"cloud.google.com/go/firestore"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

func setRoleCommand(
	ctx context.Context,
	logger *zap.SugaredLogger,
	database *firestore.Client,
	s *discordgo.Session,
	m *discordgo.MessageCreate,
) {
	// Regex match !premint-set-role <API Key>
	re := regexp.MustCompile(`^!premint-set-role (.*)$`)
	match := re.FindStringSubmatch(m.Content)

	if len(match) != 2 {
		return
	}

	// TODO: Set the Premint role
	// p := getConfig(ctx, logger, database, m.GuildID)
	// p.doc.Ref.Update(ctx, []firestore.Update{{Path: "premint-role", Value: match[1]}})
	// p.doc.Ref.Update(ctx, []firestore.Update{{Path: "premint-role-id", Value: match[1]}})
	// s.ChannelMessageSend(m.ChannelID, "✅ Premint role updated")
	s.ChannelMessageSend(m.ChannelID, "Coming soon!")
}
