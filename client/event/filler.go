package event

// func init() {
// 	List["filler"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
// 		switch i.Type {
// 		case discordgo.InteractionApplicationCommand:
// 			data := i.ApplicationCommandData()
// 			var key string
// 			if len(data.Options) > 0 {
// 				key = data.Options[0].StringValue()
// 			}

// 			list := []string{}

// 			if key != "" {
// 				room, ok := RoomList[key]
// 				if !ok {
// 					discord.EmbedError(s, i, discord.EmbedErrRoom404)
// 					return
// 				}

// 				for _, fillers := range room.Schedule {
// 					for _, f := range fillers {

// 					}
// 				}
// 			}

// 			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
// 				Type: discordgo.InteractionResponseChannelMessageWithSource,
// 				Data: &discordgo.InteractionResponseData{
// 					Embeds: []*discordgo.MessageEmbed{
// 						{
// 							Title:     fmt.Sprint(s.State.User.String(), " - Stats"),
// 							Color:     discord.EmbedColor,
// 							Timestamp: discord.EmbedTimestamp,
// 							Footer:    discord.EmbedFooter(s),
// 							Fields: []*discordgo.MessageEmbedField{
// 								{
// 									Name:  "In Guilds",
// 									Value: discord.FieldStyle(len(s.State.Guilds)),
// 								},
// 								{
// 									Name:  "Guild List",
// 									Value: discord.FieldStyle(strings.Join(guildList, ", ")),
// 								},
// 							},
// 						},
// 					},
// 					Flags: discordgo.MessageFlagsEphemeral,
// 				},
// 			})
// 			if err != nil {
// 				ErrExit(err)
// 			}

// 		// Autocomplete
// 		case discordgo.InteractionApplicationCommandAutocomplete:
// 			roomAutocomplete(s, i)
// 		}
// 	}
// }
