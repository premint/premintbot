package bot

import (
	"context"
	"regexp"

	"cloud.google.com/go/firestore"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

func setPremintCommand(
	ctx context.Context,
	logger *zap.SugaredLogger,
	database *firestore.Client,
	s *discordgo.Session,
	m *discordgo.MessageCreate,
) {
	p := getConfig(ctx, logger, database, m.GuildID)

	// Regex match !set-premint <API Key>
	re := regexp.MustCompile(`^!set-premint (.*)$`)
	match := re.FindStringSubmatch(m.Content)

	if len(match) != 2 {
		return
	}

	// Set the Premint API Key
	p.doc.Ref.Update(ctx, []firestore.Update{{Path: "premint-api-key", Value: match[1]}})
	s.ChannelMessageSend(m.ChannelID, "✅ Premint API Key updated")
}
