# TTRPG Inventory Tracker Discord Bot
A Discord Bot that helps organize your character's and/or party's inventory in a TTRPG.

## Motivation
In trying to find a good solution for keeping track of inventory for TTRPGs, I found that most options did not fit my use case.
I wanted an option that made recording and retrieving an inventory as seamless as possible, as well as allowing for listing items by user defined categories.
InventoryTracker achieves these features with the following:
- Being a discord bot makes recording and retrieving seamless since you do not have to tab over to different documents.
- Each item can be registered with a custom category to be later used in conditional item listings.

## Quick Start

### 1. Invite the bot to your discord server
If you are not the owner of the server, you need to first have permission to invite bots.

Once you do, go to the [install link](https://discord.com/oauth2/authorize?client_id=1422779683954495529)

### 2. Register your server and any campaigns/games
While in your server, enter this command to register your server: `!register server`

Then enter this command to register a game to your server: `!register game 'game_name'`

The name must only contain alphanumeric characters.

### 3. List Games
`!listGames`

Use this command to get the names of all of the games current registered to the server.

### 4. Add Players/Team Aliases
`!addPlayer 'game_name' 'player_name'`

The player's name must only contain alphanumeric characters.

### 5. List Players
`!listPlayers 'game_name'`

Use this command to get the names of all of the player's in a game. Also lists player ID.

Each command will give a feedback message on if it was successful or not and the reason if it was not successful. If successful, it will print the relevant information for that entry, such as name and ID.

### 6. Add items to player inventory
```
Add item to inventory using game and player names:
!addItem 'game_name' 'player_name' 'item_name' 'quantity'

Add item to inventory using player ID:
!addItem 'player_id' 'item_name' 'quantity'
```

### 7. List a player's inventory
```
List using game and player names:
!listInventory 'player_name' 'game_name'

List using player ID:
!listInventory 'player_id'
```

## Usage

The following are the commands you are able to run outside of the ones outlined in the Quick Start section:
- !register item 'item_name' 'game_name': registers an item to a game
- !register itemDesc 'item_name' 'game_name' 'description': registers an item with a description to a game
- !register itemCat 'item_name' 'game_name' 'category': registers an item with a category to a game
- !register itemFull 'item_name' 'game_name' 'category' 'description': registers an item with both a category and description to a game
- !updateItem description 'item_name' 'game_name' 'description': updates an existing item's description for a game
- !updateItem category 'item_name' 'game_name' 'category': updates an existing item's category for a game
- !listItem 'item_name' 'game_name': lists an item along with its description and category for a game
- !listItems 'category' 'game_name': lists all items in a category with their descriptions for a game
- !listInvByCat 'player_id' 'category': lists all items in a player's inventory with a certain category
- !listInvByCat 'player_name' 'game_name' 'category': same as above, but using player_name and game_name instead of player ID

## Contributing

The project can be tested by inviting the bot to your discord server and running its commands.

Additionally, any feature suggestions and bug fixes can be reported in issues.
