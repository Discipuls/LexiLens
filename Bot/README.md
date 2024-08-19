
# greenton-telegram-user-bot
 I'm an official GreenTON bot. Here you can plant trees at any point of the world for crypto and view your forest. Just open our Mini App and connect your TON wallet.
 
## Commands
|  Command  | Action  |
| ------------ | ------------ |
| /start  | greentings  |
| /privicy  | privicy policy  |
## Inline Query
You can use demonstration tree gift function by typing `@YourButUsername` in any chat and then select a tree for gift. After you or smb else will be able to claim the tree by clicking on "✅ Retrieve" button
## Buttons
|  Command  | Action  |
| ------------ | ------------ |
| ⚙ Settings | Open settings (currently notifications)  |
| 🌳 Plant tree | Open mini-app  |
|🔙 Cancel | Cancel action and open start message |
| Turn on notifications |   |
| Turn off notifications |   | |
## Develop project locally
### Requirements
- The bot requires environment:
```
BOT_TOKEN = ""

MINI_APP_URL = ""

DB_HOST = ""

DB_USER = ""

DB_PASSWORD = ""

DB_NAME = ""

DB_PORT = ""
```
- Or you can setup all this environment in config.json (unable for `Vercel`):
`config.json` inside `/` folder with strucure as below:
```
{
	"bot": {
		"token": "",
	},
	"mini-app": {
		"url": ""
	},
	"database": {
		"host": "",
		"user": "",
		"password": "",
		"database-name": "",
		"port": ""
	},
}
```

### Sturt up
#### Local
1. `go build`
2. `./greenton-telegram-user-bot`  or `greenton-telegram-user-bot.exe`
####  Vercel
1. Setup environment variables in your vercel project
2. Run project from github or locally with `vercel` or `vercel --prod`
3. Current version requires to setup Webhook manually:
	- `https://api.telegram.org/bot{TOKEN}/setWebhook?url={{DOMAIN}/{TOKEN}`

 
