# NotifyMe

Execute command in bash and wait to send notification via slack. It is good to save time for custom jobs when take much
time to complete and have news when is finished and profit the time.

![Slack Message](https://github.com/swapbyt3s/NotifyMe/raw/master/assents/slack.png)

## Install

Paste that at a Terminal prompt:

```bash
bash < <(curl -s https://raw.githubusercontent.com/swapbyt3s/notifyme/master/install.sh)
```

## Configure

Define this environment variables on bash before execute command.

```bash
export NOTIFYME_SLACK_TOKEN=x35QCtUUQ*B376M2D8F.JntD801gqXwOMTYuZTdGhNQ0
export NOTIFYME_SLACK_CHANNEL=alerts
```

## Usage

Is very easy to use:

```bash
./nm "mysqldump --login-path=local foo | gzip > backup.tar.gz"
```

And work on something else until you wait for notification in slack. Maybe to not lost execution, run this command
into `tmux` or `screen`.
