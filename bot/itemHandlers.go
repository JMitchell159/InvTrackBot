package bot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/JMitchell159/InvTrackBot/internal/database"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

func (st *state) addItem(s *discordgo.Session, e *discordgo.MessageCreate, cmdArgs []string) {
	/* Syntax:
	!addItem <game_name> <player_name> <item_name> <quantity>
	OR
	!addItem <player_id> <item_name> <quantity>
	*/
	if len(cmdArgs) < 3 {
		_, err := s.ChannelMessageSend(e.ChannelID, "The addItem command needs to be supplied with either 3 or 4 arguments. Look in the help section for info on addItem usage.")
		if err != nil {
			fmt.Println("Failed sending required arguments response:", err)
		}
		return
	}
	if len(cmdArgs) == 3 {
		player_id, err := uuid.Parse(cmdArgs[0])
		if err != nil {
			_, err = s.ChannelMessageSend(e.ChannelID, "Failed to parse uuid.")
			if err != nil {
				fmt.Println("Failed to send failed uuid parse response:", err)
			}
			return
		}
		quant, err := strconv.Atoi(cmdArgs[2])
		if err != nil {
			_, err = s.ChannelMessageSend(e.ChannelID, "Failed to parse quantity argument.")
			if err != nil {
				fmt.Println("Failed sending failed quantity parse response:", err)
			}
			return
		}
		_, err = st.db.GetItem(context.Background(), cmdArgs[1])
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			_, err = s.ChannelMessageSend(e.ChannelID, "Failed to fetch item.")
			if err != nil {
				fmt.Println("Failed sending failed item fetch response:", err)
			}
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			_, err = st.db.CreateItem(context.Background(), database.CreateItemParams{
				Name:      cmdArgs[1],
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
			if err != nil {
				_, err = s.ChannelMessageSend(e.ChannelID, "Failed to create item.")
				if err != nil {
					fmt.Println("Failed sending failed item create response:", err)
				}
				return
			}
			lineItem, err := st.db.AddLineItem(context.Background(), database.AddLineItemParams{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Quantity:  int32(quant),
				OwnerID:   player_id,
			})
			if err != nil {
				_, err = s.ChannelMessageSend(e.ChannelID, "Failed to add Inventory entry.")
				if err != nil {
					fmt.Println("Failed to send failed inventory add response:", err)
				}
				return
			}
			msg := ""
			if lineItem.Quantity == 1 {
				msg = fmt.Sprintf("%s has 1 %s in their inventory.", lineItem.OwnerName, lineItem.ItemName)
			}
			if lineItem.Quantity > 1 {
				msg = fmt.Sprintf("%s has %d %ss in their inventory.", lineItem.OwnerName, lineItem.Quantity, lineItem.ItemName)
			}
			_, err = s.ChannelMessageSend(e.ChannelID, msg)
			if err != nil {
				fmt.Println("Failed to send addItem response:", err)
			}
			return
		}
		_, err = st.db.GetPlayer(context.Background(), player_id)
		if errors.Is(err, sql.ErrNoRows) {
			_, err = s.ChannelMessageSend(e.ChannelID, fmt.Sprintf("Player_id %s does not exist.", player_id.String()))
			if err != nil {
				fmt.Println("Failed to send invalid player_id response:", err)
			}
			return
		}
		_, err = st.db.GetLineItemByItemAndOwner(context.Background(), database.GetLineItemByItemAndOwnerParams{
			OwnerID:  player_id,
			ItemName: cmdArgs[1],
		})
		if errors.Is(err, sql.ErrNoRows) {
			lineItem, err := st.db.AddLineItem(context.Background(), database.AddLineItemParams{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Quantity:  int32(quant),
				OwnerID:   player_id,
			})
			if err != nil {
				_, err = s.ChannelMessageSend(e.ChannelID, "Failed to add Inventory entry.")
				if err != nil {
					fmt.Println("Failed to send failed inventory add response:", err)
				}
				return
			}
			msg := ""
			if lineItem.Quantity == 1 {
				msg = fmt.Sprintf("%s has 1 %s in their inventory.", lineItem.OwnerName, lineItem.ItemName)
			}
			if lineItem.Quantity > 1 {
				msg = fmt.Sprintf("%s has %d %ss in their inventory.", lineItem.OwnerName, lineItem.Quantity, lineItem.ItemName)
			}
			_, err = s.ChannelMessageSend(e.ChannelID, msg)
			if err != nil {
				fmt.Println("Failed to send addItem response:", err)
			}
			return
		}
		err = st.db.UpdateLineItem(context.Background(), database.UpdateLineItemParams{
			Quantity:  int32(quant),
			UpdatedAt: time.Now(),
			OwnerID:   player_id,
			ItemName:  cmdArgs[1],
		})
		if err != nil {
			_, err = s.ChannelMessageSend(e.ChannelID, "Failed to update inventory.")
			if err != nil {
				fmt.Println("Failed to send failed inventory update response:", err)
			}
			return
		}
		player_name, _ := st.db.GetPlayer(context.Background(), player_id)
		msg := ""
		if quant == 1 {
			msg = fmt.Sprintf("%s added 1 %s to their inventory.", player_name, cmdArgs[1])
		} else if quant > 1 {
			msg = fmt.Sprintf("%s added %d %ss to their inventory.", player_name, quant, cmdArgs[1])
		}
		_, err = s.ChannelMessageSend(e.ChannelID, msg)
		if err != nil {
			fmt.Println("Failed to send addItem response:", err)
		}
		return
	}
	quant, err := strconv.Atoi(cmdArgs[3])
	if err != nil {
		_, err = s.ChannelMessageSend(e.ChannelID, "Failed to parse quantity.")
		if err != nil {
			fmt.Println("Failed to send failed quantity parse response:", err)
		}
		return
	}
	game, err := st.db.GetGameByName(context.Background(), database.GetGameByNameParams{
		Name:     cmdArgs[0],
		ServerID: e.GuildID,
	})
	if errors.Is(err, sql.ErrNoRows) {
		_, err = s.ChannelMessageSend(e.ChannelID, fmt.Sprintf("Game %s does not exist in this server.", cmdArgs[0]))
		if err != nil {
			fmt.Println("Failed to send invalid game response:", err)
		}
		return
	}
	if err != nil {
		_, err = s.ChannelMessageSend(e.ChannelID, "Something went wrong while fetching game.")
		if err != nil {
			fmt.Println("Failed to send failed game fetch response:", err)
		}
		return
	}
	player, err := st.db.GetPlayerByName(context.Background(), database.GetPlayerByNameParams{
		Name:   cmdArgs[1],
		GameID: game.ID,
	})
	if errors.Is(err, sql.ErrNoRows) {
		_, err = s.ChannelMessageSend(e.ChannelID, fmt.Sprintf("player %s does not exist in this game.", cmdArgs[1]))
		if err != nil {
			fmt.Println("Failed to send invalid player response:", err)
		}
		return
	}
	if err != nil {
		_, err = s.ChannelMessageSend(e.ChannelID, "Something went wrong while fetching player.")
		if err != nil {
			fmt.Println("Failed to send failed player fetch response:", err)
		}
		return
	}
	_, err = st.db.GetItem(context.Background(), cmdArgs[2])
	if errors.Is(err, sql.ErrNoRows) {
		_, err = st.db.CreateItem(context.Background(), database.CreateItemParams{
			Name:      cmdArgs[2],
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		if err != nil {
			_, err = s.ChannelMessageSend(e.ChannelID, "Failed to create item.")
			if err != nil {
				fmt.Println("Failed sending failed item create response:", err)
			}
			return
		}
		lineItem, err := st.db.AddLineItem(context.Background(), database.AddLineItemParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Quantity:  int32(quant),
			OwnerID:   player.ID,
		})
		if err != nil {
			_, err = s.ChannelMessageSend(e.ChannelID, "Failed to add Inventory entry.")
			if err != nil {
				fmt.Println("Failed to send failed inventory add response:", err)
			}
			return
		}
		msg := ""
		if lineItem.Quantity == 1 {
			msg = fmt.Sprintf("%s has 1 %s in their inventory.", lineItem.OwnerName, lineItem.ItemName)
		}
		if lineItem.Quantity > 1 {
			msg = fmt.Sprintf("%s has %d %ss in their inventory.", lineItem.OwnerName, lineItem.Quantity, lineItem.ItemName)
		}
		_, err = s.ChannelMessageSend(e.ChannelID, msg)
		if err != nil {
			fmt.Println("Failed to send addItem response:", err)
		}
		return
	}
	_, err = st.db.GetLineItemByItemAndOwner(context.Background(), database.GetLineItemByItemAndOwnerParams{
		OwnerID:  player.ID,
		ItemName: cmdArgs[2],
	})
	if errors.Is(err, sql.ErrNoRows) {
		lineItem, err := st.db.AddLineItem(context.Background(), database.AddLineItemParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Quantity:  int32(quant),
			OwnerID:   player.ID,
			ItemName:  cmdArgs[2],
		})
		if err != nil {
			_, err = s.ChannelMessageSend(e.ChannelID, "Something went wrong while adding inventory entry.")
			if err != nil {
				fmt.Println("Failed to send failed inventory addition response:", err)
			}
			return
		}
		msg := ""
		if lineItem.Quantity == 1 {
			msg = fmt.Sprintf("%s has 1 %s in their inventory.", lineItem.OwnerName, lineItem.ItemName)
		}
		if lineItem.Quantity > 1 {
			msg = fmt.Sprintf("%s has %d %ss in their inventory.", lineItem.OwnerName, lineItem.Quantity, lineItem.ItemName)
		}
		_, err = s.ChannelMessageSend(e.ChannelID, msg)
		if err != nil {
			fmt.Println("Failed to send addItem response:", err)
		}
		return
	}
	err = st.db.UpdateLineItem(context.Background(), database.UpdateLineItemParams{
		Quantity:  int32(quant),
		UpdatedAt: time.Now(),
		OwnerID:   player.ID,
		ItemName:  cmdArgs[2],
	})
	if err != nil {
		_, err = s.ChannelMessageSend(e.ChannelID, "Failed to update inventory.")
		if err != nil {
			fmt.Println("Failed to send failed inventory update response:", err)
		}
		return
	}
	msg := ""
	if quant == 1 {
		msg = fmt.Sprintf("%s added 1 %s to their inventory.", cmdArgs[1], cmdArgs[2])
	} else if quant > 1 {
		msg = fmt.Sprintf("%s added %d %ss to their inventory.", cmdArgs[1], quant, cmdArgs[2])
	}
	_, err = s.ChannelMessageSend(e.ChannelID, msg)
	if err != nil {
		fmt.Println("Failed to send addItem response:", err)
	}
}
