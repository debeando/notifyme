# NotifyMe

Execute command in bash and wait to send notification via slack.

![Slack Message](https://github.com/swapbyt3s/NotifyMe/raw/master/assents/slack.png)


## Setup

Define this environment variables on bash before execute command.

```bash
export NOTIFYME_SLACK_TOKEN=x35QCtUUQ*B376M2D8F.JntD801gqXwOMTYuZTdGhNQ0
export NOTIFYME_SLACK_CHANNEL=alerts
```

## Usage

Is very easy to use:

```bash
./notifyme "ls -lah"
```

And work on something else until you wait for notification in slack. Maybe to not lost execution, run this command into `tmux` or `screen`.
