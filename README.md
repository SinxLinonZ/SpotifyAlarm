# SpotifyAlarm
Schedule alarms to play Spotify on any online device.

# Usage
Warning: This is a work in progress, use at your own risk. Currently only supports Windows.
1. Create an app on the [Spotify Developer Dashboard](https://developer.spotify.com/dashboard/applications)
2. Set the redirect URI to `http://127.0.0.1:13333/callback`
3. Copy `config.example.yaml` to `config.yaml` and fill in the redirect URI, client ID and secret
4. Run `SpotifyAlarm.exe`
5. For the first time, an authentication window will open. Log in to Spotify and allow the app to access your account.
6. An icon should be in your system tray. Right click it to open UI.


# Development
TODO:
- [ ] Edit existing alarms
- [ ] Delete existing alarms
- [ ] Support for playlists
- [ ] Make an icon

# Thanks
- [Spotify Web API](https://developer.spotify.com/web-api/)
- [Spotify Go wrapper](https://github.com/zmb3/spotify)
- [systray](https://github.com/getlantern/systray)