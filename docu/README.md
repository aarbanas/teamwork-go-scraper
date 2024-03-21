# How to setup for local usage

## Download binary

Go to the [releases page](https://github.com/aarbanas/teamwork-go-scraper/releases) and download binary specific four your operating system.

1. Windows users - most of todays PC's are using 64-bit processor, which means that `teamwork-go-scraper_0.0.0_windows_amd64.tar.gz` is the correct version.
2. Mac users just need to take care if they are using M-series of processors than `teamwork-go-scraper_0.0.0_darwin_arm64.tar.gz` is the one, but if you are using old Intel processor, than download the `teamwork-go-scraper_0.0.0_darwin_amd64.tar.gz`
3. Linux users - same as the Windows users

## Create config.json

In the directory where your binary is downloaded (it is usually `Downloads` folder) create `config.json` file.

For Linux/Mac and Windows WSL users (open terminal):

```shell
$ cd Downloads
$ nano config.json
# or for vim users
$ vi config.json
```

For Windows users not using WSL (open cmd):

```powershell
cd Downloads
echo. > config.json
notepad config.json
```

After you open `config.json` in your new editor, copy the content from [config.json-example](https://github.com/aarbanas/teamwork-go-scraper/blob/main/config.json-example) and paste it in the editor.

### Enter required data to config.json

1. `userId` - In Teamwork open your profile
   1. Click on profile image (bottom left corner)
   2. View profile
   3. Copy the userId from the URL: https://teamwork/app/people/COPY_THIS_VALUE/time
2. `apiKey` - In Teamwork open edit your details
   1. Click on profile image (bottom left corner)
   2. Edit my details
   3. API & Mobile
   4. Your API Token (Generate new or Show your token if you already have one)
3. `url` - The base URL from your teamwork account
4. `userIds` - enter user ids from all team members you whish to check hours periodically

> Remember to save changes made and you are ready to go.
