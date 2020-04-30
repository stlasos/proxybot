# Proxybot

Telegram bot which checks if urls can be used as a proxy.

You can fill in ip and port lists and proxybot will check all their possible combinations for the ability to proxy through these addresses.

The following config is required:

```json
{
  "tg_api_url": "",
  "tg_token": "",
  "proxy_connect_url": ""
}
```

You can create a config file named `proxybot.json` (which is used by default) or specify the config location by passing an argument to the command.

`tg_api_url` - Url for the telegram API (if empty https://api.telegram.org will be used)

`tg_token` - Your bot's token

`proxy_connect_url` - Url to connect via proxy (200 status code should be returned)