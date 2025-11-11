package bot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
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
		sendMessage(s, e.ChannelID, "The addItem command needs to be supplied with either 3 or 4 arguments. Look in the help section for info on addItem usage.", "Failed sending required arguments response:")
		return
	}
	if len(cmdArgs) == 3 {
		player_id, err := uuid.Parse(cmdArgs[0])
		if err != nil {
			sendMessage(s, e.ChannelID, "Failed to parse uuid.", "Failed sending failed uuid parse response:")
			return
		}
		player, err := st.db.GetPlayer(context.Background(), player_id)
		if errors.Is(err, sql.ErrNoRows) {
			sendMessage(s, e.ChannelID, "Specified Player ID does not exist.", "failed sending invalid Player ID response:")
			return
		}
		quant, err := strconv.Atoi(cmdArgs[2])
		if err != nil {
			sendMessage(s, e.ChannelID, "Failed to parse quantity argument.", "Failed sending failed quantity parse response:")
			return
		}
		itemWithHash := fmt.Sprintf("%s[%s]", cmdArgs[1], player.GameID.String())
		_, err = st.db.GetItem(context.Background(), itemWithHash)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			sendMessage(s, e.ChannelID, "Failed to fetch item.", "Failed sending failed item fetch response:")
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			_, err = st.db.CreateItem(context.Background(), database.CreateItemParams{
				Name:      itemWithHash,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
			if err != nil {
				sendMessage(s, e.ChannelID, "Failed to create item.", "Failed sending failed item create response:")
				return
			}
			lineItem, err := st.db.AddLineItem(context.Background(), database.AddLineItemParams{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Quantity:  int32(quant),
				OwnerID:   player_id,
				ItemName:  itemWithHash,
			})
			if err != nil {
				sendMessage(s, e.ChannelID, "Failed to add Inventory entry.", "Failed to send failed inventory add response:")
				return
			}
			msg := ""
			if lineItem.Quantity == 1 {
				msg = fmt.Sprintf("%s has 1 %s in their inventory.", lineItem.OwnerName, lineItem.ItemName)
			}
			if lineItem.Quantity > 1 {
				msg = fmt.Sprintf("%s has %d %ss in their inventory.", lineItem.OwnerName, lineItem.Quantity, lineItem.ItemName)
			}
			sendMessage(s, e.ChannelID, msg, "Failed to send addItem response:")
			return
		}
		_, err = st.db.GetLineItemByItemAndOwner(context.Background(), database.GetLineItemByItemAndOwnerParams{
			OwnerID:  player_id,
			ItemName: itemWithHash,
		})
		if errors.Is(err, sql.ErrNoRows) {
			lineItem, err := st.db.AddLineItem(context.Background(), database.AddLineItemParams{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Quantity:  int32(quant),
				OwnerID:   player_id,
				ItemName:  itemWithHash,
			})
			if err != nil {
				sendMessage(s, e.ChannelID, "Failed to add Inventory entry.", "Failed to send failed inventory add response:")
				return
			}
			msg := ""
			split := strings.Split(lineItem.ItemName, "[")
			if lineItem.Quantity == 1 {
				msg = fmt.Sprintf("%s has 1 %s in their inventory.", lineItem.OwnerName, split[0])
			}
			if lineItem.Quantity > 1 {
				msg = fmt.Sprintf("%s has %d %ss in their inventory.", lineItem.OwnerName, lineItem.Quantity, split[0])
			}
			sendMessage(s, e.ChannelID, msg, "Failed to send addItem response:")
			return
		}
		err = st.db.UpdateLineItemWID(context.Background(), database.UpdateLineItemWIDParams{
			Quantity:  int32(quant),
			UpdatedAt: time.Now(),
			OwnerID:   player_id,
			ItemName:  itemWithHash,
		})
		if err != nil {
			sendMessage(s, e.ChannelID, "Failed to update inventory.", "Failed to send failed inventory update response:")
			return
		}
		player, _ = st.db.GetPlayer(context.Background(), player_id)
		msg := ""
		if quant == 1 {
			msg = fmt.Sprintf("%s added 1 %s to their inventory.", player.Name, cmdArgs[1])
		} else if quant > 1 {
			msg = fmt.Sprintf("%s added %d %ss to their inventory.", player.Name, quant, cmdArgs[1])
		}
		sendMessage(s, e.ChannelID, msg, "Failed to send addItem response:")
		return
	}
	quant, err := strconv.Atoi(cmdArgs[3])
	if err != nil {
		sendMessage(s, e.ChannelID, "Failed to parse quantity.", "Failed to send failed quantity parse response:")
		return
	}
	game, err := st.db.GetGameByName(context.Background(), database.GetGameByNameParams{
		Name:     cmdArgs[0],
		ServerID: e.GuildID,
	})
	if errors.Is(err, sql.ErrNoRows) {
		sendMessage(s, e.ChannelID, fmt.Sprintf("Game %s does not exist in this server.", cmdArgs[0]), "Failed to send invalid game response:")
		return
	}
	if err != nil {
		sendMessage(s, e.ChannelID, "Something went wrong while fetching game.", "Failed to send fialed game fetch response:")
		return
	}
	player, err := st.db.GetPlayerByName(context.Background(), database.GetPlayerByNameParams{
		Name:   cmdArgs[1],
		GameID: game.ID,
	})
	if errors.Is(err, sql.ErrNoRows) {
		sendMessage(s, e.ChannelID, fmt.Sprintf("player %s does not exist in this game.", cmdArgs[1]), "Failed to send invalid player response:")
		return
	}
	if err != nil {
		sendMessage(s, e.ChannelID, "Something went wrong while fetching player.", "Failed to send failed player fetch response:")
		return
	}
	itemWithHash := fmt.Sprintf("%s[%s]", cmdArgs[2], game.ID.String())
	_, err = st.db.GetItem(context.Background(), itemWithHash)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = st.db.CreateItem(context.Background(), database.CreateItemParams{
			Name:      itemWithHash,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		if err != nil {
			sendMessage(s, e.ChannelID, "Failed to create item.", "Failed sending failed item create response:")
			return
		}
		lineItem, err := st.db.AddLineItem(context.Background(), database.AddLineItemParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Quantity:  int32(quant),
			OwnerID:   player.ID,
			ItemName:  itemWithHash,
		})
		if err != nil {
			sendMessage(s, e.ChannelID, "Failed to add Inventory entry.", "Failed to send failed inventory add response:")
			return
		}
		msg := ""
		split := strings.Split(lineItem.ItemName, "[")
		if lineItem.Quantity == 1 {
			msg = fmt.Sprintf("%s has 1 %s in their inventory.", lineItem.OwnerName, split[0])
		}
		if lineItem.Quantity > 1 {
			msg = fmt.Sprintf("%s has %d %ss in their inventory.", lineItem.OwnerName, lineItem.Quantity, split[0])
		}
		sendMessage(s, e.ChannelID, msg, "Failed to send addItem response:")
		return
	}
	_, err = st.db.GetLineItemByItemAndOwner(context.Background(), database.GetLineItemByItemAndOwnerParams{
		OwnerID:  player.ID,
		ItemName: itemWithHash,
	})
	if errors.Is(err, sql.ErrNoRows) {
		lineItem, err := st.db.AddLineItem(context.Background(), database.AddLineItemParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Quantity:  int32(quant),
			OwnerID:   player.ID,
			ItemName:  itemWithHash,
		})
		if err != nil {
			sendMessage(s, e.ChannelID, "Something went wrong while adding inventory entry.", "Failed to send failed inventory addition response:")
			return
		}
		msg := ""
		split := strings.Split(lineItem.ItemName, "[")
		if lineItem.Quantity == 1 {
			msg = fmt.Sprintf("%s has 1 %s in their inventory.", lineItem.OwnerName, split[0])
		}
		if lineItem.Quantity > 1 {
			msg = fmt.Sprintf("%s has %d %ss in their inventory.", lineItem.OwnerName, lineItem.Quantity, split[0])
		}
		sendMessage(s, e.ChannelID, msg, "Failed to send addItem repsonse:")
		return
	}
	err = st.db.UpdateLineItemWID(context.Background(), database.UpdateLineItemWIDParams{
		Quantity:  int32(quant),
		UpdatedAt: time.Now(),
		OwnerID:   player.ID,
		ItemName:  itemWithHash,
	})
	if err != nil {
		sendMessage(s, e.ChannelID, "Failed to update inventory.", "Failed to send failed inventory update response:")
		return
	}
	msg := ""
	if quant == 1 {
		msg = fmt.Sprintf("%s added 1 %s to their inventory.", cmdArgs[1], cmdArgs[2])
	} else if quant > 1 {
		msg = fmt.Sprintf("%s added %d %ss to their inventory.", cmdArgs[1], quant, cmdArgs[2])
	}
	sendMessage(s, e.ChannelID, msg, "Failed to send addItem response:")
}

func (st *state) updateItem(s *discordgo.Session, e *discordgo.MessageCreate, cmdArgs []string) {
	if len(cmdArgs) < 4 {
		sendMessage(s, e.ChannelID, "The updateItem command must have 4 arguments. The syntax is as follows: '!updateItem description|category <item_name> <game_name> <description|category>.'", "Failed to send required arguments response:")
		return
	}

	game, err := st.db.GetGameByName(context.Background(), database.GetGameByNameParams{
		Name:     cmdArgs[2],
		ServerID: e.GuildID,
	})
	if errors.Is(err, sql.ErrNoRows) {
		sendMessage(s, e.ChannelID, "Specified game does not exist in this server.", "Failed to send invalid Game response:")
		return
	}

	itemWithHash := fmt.Sprintf("%s[%s]", cmdArgs[1], game.ID.String())
	_, err = st.db.GetItem(context.Background(), itemWithHash)
	if errors.Is(err, sql.ErrNoRows) {
		sendMessage(s, e.ChannelID, "That item does not exist. Please register it first.", "Failed to send invalid item response:")
		return
	}

	/*Syntax:
	!updateItem description <item_name> <game_name> <description>*/
	if strings.ToLower(cmdArgs[0]) == "description" {
		err = st.db.UpdateDesc(context.Background(), database.UpdateDescParams{
			Description: sql.NullString{
				String: cmdArgs[2],
				Valid:  true,
			},
			UpdatedAt: time.Now(),
			Name:      itemWithHash,
		})
		if err != nil {
			sendMessage(s, e.ChannelID, "Failed to update description.", "Failed to send failed description update response:")
			return
		}
		sendMessage(s, e.ChannelID, fmt.Sprintf("Updated %s's description to %s.", cmdArgs[1], cmdArgs[3]), "Failed to send description update response:")
		return
	}

	/*Syntax:
	!updateItem category <item_name> <game_name> <category>*/
	if strings.ToLower(cmdArgs[0]) == "category" {
		err = st.db.UpdateCat(context.Background(), database.UpdateCatParams{
			Category: sql.NullString{
				String: cmdArgs[2],
				Valid:  true,
			},
			UpdatedAt: time.Now(),
			Name:      itemWithHash,
		})
		if err != nil {
			sendMessage(s, e.ChannelID, "Failed to update category.", "Failed to send failed category update response:")
			return
		}
		sendMessage(s, e.ChannelID, fmt.Sprintf("Updated %s's category to %s.", cmdArgs[1], cmdArgs[3]), "Failed to send category update response:")
		return
	}

	sendMessage(s, e.ChannelID, fmt.Sprintf("Unknown command updateItem %s.", cmdArgs[0]), "Failed sending Unknown Command response:")
}

func (st *state) listItem(s *discordgo.Session, e *discordgo.MessageCreate, cmdArgs []string) {
	/*Syntax:
	!listItem <item_name> <game_name>*/
	game, err := st.db.GetGameByName(context.Background(), database.GetGameByNameParams{
		Name:     cmdArgs[2],
		ServerID: e.GuildID,
	})
	if errors.Is(err, sql.ErrNoRows) {
		sendMessage(s, e.ChannelID, "Specified game does not exist in this server.", "Failed to send invalid Game response:")
		return
	}
	itemWithHash := fmt.Sprintf("%s[%s]", cmdArgs[0], game.ID.String())
	item, err := st.db.GetItem(context.Background(), itemWithHash)
	if errors.Is(err, sql.ErrNoRows) {
		sendMessage(s, e.ChannelID, "That item does not exist.", "Failed to send invalid item response:")
		return
	}
	if err != nil {
		sendMessage(s, e.ChannelID, "Something went wrong while fetching item.", "Failed to send failed item fetch response:")
		return
	}
	split := strings.Split(item.Name, "[")
	msg := fmt.Sprintf("Name:\t\t\t%s\n", split[0])
	if item.Description.Valid {
		msg = fmt.Sprintf("%sDescription:\t%s\n", msg, item.Description.String)
	} else {
		msg = fmt.Sprintf("%sDescription:\tNo Description\n", msg)
	}
	if item.Category.Valid {
		msg = fmt.Sprintf("%sCategory:\t\t%s", msg, item.Category.String)
	} else {
		msg = fmt.Sprintf("%sCategory:\t\tNo Category", msg)
	}
	sendMessage(s, e.ChannelID, msg, "Failed to send listItem response:")
}

func (st *state) listItems(s *discordgo.Session, e *discordgo.MessageCreate, cmdArgs []string) {
	/*Syntax:
	!listItems <category> <game_name>*/
	game, err := st.db.GetGameByName(context.Background(), database.GetGameByNameParams{
		Name:     cmdArgs[2],
		ServerID: e.GuildID,
	})
	if errors.Is(err, sql.ErrNoRows) {
		sendMessage(s, e.ChannelID, "Specified game does not exist in this server.", "Failed to send invalid Game response:")
		return
	}
	items, err := st.db.GetItemsByCategory(context.Background(), sql.NullString{
		Valid:  true,
		String: cmdArgs[0],
	})
	if err != nil {
		sendMessage(s, e.ChannelID, "Something went wrong while fetching item.", "Failed to send failed item fetch response:")
		return
	}
	if len(items) == 0 {
		sendMessage(s, e.ChannelID, "There are no items in that category.", "Failed to send invalid category response:")
		return
	}
	msg := fmt.Sprintf("Items in %s category:\n", cmdArgs[0])
	for _, item := range items {
		if strings.Contains(item.Name, game.ID.String()) {
			split := strings.Split(item.Name, "[")
			if item.Description.Valid {
				msg = fmt.Sprintf("%s- %s: %s\n", msg, split[0], item.Description.String)
			} else {
				msg = fmt.Sprintf("%s- %s\n", msg, split[0])
			}
		}
	}
	sendMessage(s, e.ChannelID, msg, "Failed to send listItems response:")
}

func (st *state) listInventory(s *discordgo.Session, e *discordgo.MessageCreate, cmdArgs []string) {
	/*Syntax:
	!listInventory <player_id>
	OR
	!listInventory <player_name> <game_name>*/
	if len(cmdArgs) == 1 {
		player_id, err := uuid.Parse(cmdArgs[0])
		if err != nil {
			sendMessage(s, e.ChannelID, "Something went wrong while parsing player id string.", "Failed to send failed uuid parse response:")
			return
		}

		inventory, err := st.db.GetItemsByOwner(context.Background(), player_id)
		if err != nil {
			sendMessage(s, e.ChannelID, "Something went wrong while fetching inventory.", "Failed to send failed inventory fetch response:")
			return
		}
		if len(inventory) == 0 {
			sendMessage(s, e.ChannelID, "This inventory is empty, you should add some items to it.", "Failed to send empty inventory response:")
			return
		}

		player, err := st.db.GetPlayer(context.Background(), player_id)
		if err != nil {
			sendMessage(s, e.ChannelID, "Something went wrong while fetching player.", "Failed to send failed player fetch response:")
			return
		}

		sendMessage(s, e.ChannelID, fmt.Sprintf("%s's Inventory:", player.Name), "Failed:")
		msg := ""
		for _, item := range inventory {
			split := strings.Split(item.Name, "[")
			if item.Description.Valid {
				msg += fmt.Sprintf("%dx %s: %s\n", item.Quantity, split[0], item.Description.String)
			} else {
				msg += fmt.Sprintf("%dx %s\n", item.Quantity, split[0])
			}
			if item.Category.Valid {
				msg += fmt.Sprintf("Category: %s", item.Category.String)
			}
			sendMessage(s, e.ChannelID, msg, "Failed:")
			msg = ""
		}
		sendMessage(s, e.ChannelID, "End of Inventory", "Failed to send inventory response:")
		return
	}

	inventory, err := st.db.GetItemsByOwnerName(context.Background(), database.GetItemsByOwnerNameParams{
		Name:     cmdArgs[0],
		Name_2:   cmdArgs[1],
		ServerID: e.GuildID,
	})
	if err != nil {
		sendMessage(s, e.ChannelID, "Something went wrong while fetching inventory.", "Failed to send failed inventory fetch response:")
		return
	}
	if len(inventory) == 0 {
		sendMessage(s, e.ChannelID, "This inventory is empty, you should add some items to it.", "Failed to send empty inventory response:")
		return
	}

	sendMessage(s, e.ChannelID, fmt.Sprintf("%s's Inventory:", cmdArgs[0]), "Failed:")
	msg := ""
	for _, item := range inventory {
		split := strings.Split(item.Name, "[")
		if item.Description.Valid {
			msg += fmt.Sprintf("%dx %s: %s\n", item.Quantity, split[0], item.Description.String)
		} else {
			msg += fmt.Sprintf("%dx %s\n", item.Quantity, split[0])
		}
		if item.Category.Valid {
			msg += fmt.Sprintf("Category: %s", item.Category.String)
		}
		sendMessage(s, e.ChannelID, msg, "Failed:")
		msg = ""
	}
	sendMessage(s, e.ChannelID, "End of Inventory", "Failed to send inventory response:")
}

func (st *state) listInventoryByCat(s *discordgo.Session, e *discordgo.MessageCreate, cmdArgs []string) {
	/*Syntax:
	!listInvByCat <player_id> <category>
	OR
	!listInvByCat <player_name> <game_name> <category>*/
	if len(cmdArgs) < 2 {
		sendMessage(s, e.ChannelID, "The list by category command takes at least 2 arguments. Look in the help section for more info on usage.", "Failed to send required arguments response:")
		return
	}
	if len(cmdArgs) == 2 {
		player_id, err := uuid.Parse(cmdArgs[0])
		if err != nil {
			sendMessage(s, e.ChannelID, "Something went wrong while parsing player id string.", "Failed to send failed uuid parse response:")
			return
		}

		inventory, err := st.db.GetItemsByOwnerAndCat(context.Background(), database.GetItemsByOwnerAndCatParams{
			OwnerID: player_id,
			Category: sql.NullString{
				Valid:  true,
				String: cmdArgs[1],
			},
		})
		if err != nil {
			sendMessage(s, e.ChannelID, "Something went wrong while fetching inventory.", "Failed to send failed inventory fetch response:")
			return
		}
		if len(inventory) == 0 {
			sendMessage(s, e.ChannelID, "You have no items in that category.", "Failed to send empty list response:")
			return
		}

		player, err := st.db.GetPlayer(context.Background(), player_id)
		if err != nil {
			sendMessage(s, e.ChannelID, "Something went wrong while fetching player.", "Failed to send failed player fetch response:")
			return
		}

		sendMessage(s, e.ChannelID, fmt.Sprintf("%s's Inventory (Category: %s):", player.Name, cmdArgs[1]), "Failed:")
		msg := ""
		for _, item := range inventory {
			split := strings.Split(item.Name, "[")
			if item.Description.Valid {
				msg += fmt.Sprintf("%dx %s: %s\n", item.Quantity, split[0], item.Description.String)
			} else {
				msg += fmt.Sprintf("%dx %s\n", item.Quantity, split[0])
			}
			sendMessage(s, e.ChannelID, msg, "Failed:")
			msg = ""
		}
		sendMessage(s, e.ChannelID, "End of Inventory", "Failed to send inventory response:")
		return
	}

	inventory, err := st.db.GetItemsByOwnerNameAndCat(context.Background(), database.GetItemsByOwnerNameAndCatParams{
		Name:     cmdArgs[0],
		Name_2:   cmdArgs[1],
		ServerID: e.GuildID,
		Category: sql.NullString{
			Valid:  true,
			String: cmdArgs[2],
		},
	})
	if err != nil {
		sendMessage(s, e.ChannelID, "Something went wrong while fetching inventory.", "Failed to send failed inventory fetch response:")
		return
	}
	if len(inventory) == 0 {
		sendMessage(s, e.ChannelID, "You have no items in that category.", "Failed to send empty list response:")
		return
	}

	sendMessage(s, e.ChannelID, fmt.Sprintf("%s's Inventory (Category: %s):", cmdArgs[0], cmdArgs[2]), "Failed:")
	msg := ""
	for _, item := range inventory {
		split := strings.Split(item.Name, "[")
		if item.Description.Valid {
			msg += fmt.Sprintf("%dx %s: %s\n", item.Quantity, split[0], item.Description.String)
		} else {
			msg += fmt.Sprintf("%dx %s\n", item.Quantity, split[0])
		}
		if item.Category.Valid {
			msg += fmt.Sprintf("Category: %s", item.Category.String)
		}
		sendMessage(s, e.ChannelID, msg, "Failed:")
		msg = ""
	}
	sendMessage(s, e.ChannelID, "End of Inventory", "Failed to send inventory response:")
}
