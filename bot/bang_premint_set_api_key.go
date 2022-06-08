package bot

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/firestore"
	bq "github.com/premint/premintbot/bigquery"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

func premintSetAPIKeyCommand(
	ctx context.Context,
	logger *zap.SugaredLogger,
	database *firestore.Client,
	bqClient *bigquery.Client,
	s *discordgo.Session,
	m *discordgo.MessageCreate,
) {
	if m.Content == "!premint-set-api-key" {
		s.ChannelMessageSend(m.ChannelID, "Missing API key. Please use `!premint-set-api-key PREMINT_API_KEY` to set it. You can find your API key on the Premint website: https://www.premint.xyz/dashboard/. Click on a project, then click Edit Settings, then API.")
		bq.RecordAdminAction(bqClient, m, "set-api-key", "missing-api-key")
		return
	}

	// Regex match !premint-set-api-key PREMINT_API_KEY
	re := regexp.MustCompile(`^!premint-set-api-key (.*)$`)
	match := re.FindStringSubmatch(m.Content)

	if len(match) != 2 {
		return
	}

	p := GetConfig(ctx, logger, database, m.GuildID)

	if !isAdmin(p, m.Author) {
		s.ChannelMessageSend(m.ChannelID, "❌ You do not have the Premintbot role. Please contact a server administrator to add it to your account.")
		bq.RecordAdminAction(bqClient, m, "set-api-key", "not-admin")
		return
	}

	apiKey := match[1]

	// If apiKey contains any of the bad characters, return an error
	for _, c := range badChars {
		if strings.Contains(apiKey, c) {
			s.ChannelMessageSend(m.ChannelID, "❌ API key contains invalid characters. Please use `!premint-set-api-key PREMINT_API_KEY` to set it. You can find your API key on the Premint website: https://www.premint.xyz/dashboard/. Click on a project, then click Edit Settings, then API.")
			bq.RecordAdminAction(bqClient, m, "set-api-key", "bad-characters")
			return
		}
	}

	p.doc.Ref.Update(ctx, []firestore.Update{
		{Path: "premint-api-key", Value: apiKey},
	})

	bq.RecordAdminAction(bqClient, m, "set-api-key", "success")

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("✅ Premint API key updated: %s", match[1]))
}
