# BingGO
## About The Project
BingGo downloads and sets www.bing.com daily images to your Raspberry pi wallpaper and all major desktops.   
## Supported desktops

* Raspbian
* Windows
* macOS
* GNOME
* KDE
* Cinnamon
* Unity
* Budgie
* XFCE
* LXDE
* MATE
* Deepin
### How to Install
Download [BingGo](https://github.com/Abbas-gheydi/binggo/releases) and install it by:  
```bash
sudo install ./binggo /usr/local/bin/
```

### How to Use it:
just run it:   
wallpapers download folder is /home/pi/Pictures/binggo  
```
binggo
```
to download and set wallpapers automatically set cronjob in /etc/crontab (to download wallpapers at 10:00 for instanse)   
```
00 10	* * *	pi	binggo

```
