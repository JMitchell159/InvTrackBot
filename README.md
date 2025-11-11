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
While in your server, enter this command to register your server:
`!register server`

Then enter this command to register a game to your server:
`!register game 'game_name'`
The name must only contain alphanumeric characters.

### 3. Add Players/Team Aliases
`!addPlayer 'game_name' 'player_name'`
The player's name number only contain alphanumeric characters.

Each command will give a feedback message on if it was successful or not and the reason if it was not successful. If successful, it will print the relevant information for that entry, such as name and ID.

### 4. Add items to player inventory
```
Add item to inventory using game and player names:
!addItem 'game_name' 'player_name' 'item_name' 'quantity'

Add item to inventory using player ID:
!addItem 'player_id' 'item_name' 'quantity'
```

## Usage

## Contributing
