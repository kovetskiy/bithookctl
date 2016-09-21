# bithookctl

Manage your hooks in Bitbucket (Atlassian Stash) repository.

You should create configuration file at ~/.config/bithookctl.conf
with following syntax:

```
  user = "username"
  pass = "password"
```

If you will work with multiple hooks you can add aliases section with following
syntax:

```
  [aliases]
   pre = "com.ngs.stash.externalhooks.external-hooks:external-pre-receive-hook"
   post = "com.ngs.stash.externalhooks.external-hooks:external-post-receive-hook"
```

### Synopsis

```
Usage:
    bithookctl [options] -L
    bithookctl [options] -A <hook> <id> [(-v <var_name> <var_value>)]...
    bithookctl [options] -U <hook> <id>
    bithookctl [options] -R <hook> <id>
    bithookctl -h | --help

Options:
    -L --list       List installed hooks.
    -A --add        Add hook <hook> with <id>, compile <templates>/<hook>
                     template and use as <hook> args.
    -U --update     Update hook <hook> with <id> args.
    -R --remove     Remove hook <id> with <id>.
    -u <url>        Specify repository URL.
                     By default, it will be read from 'git remote origin'
                     output.
    -c <config>     Specify configuration file with user credentials.
                     [default: $HOME/.config/bithookctl.conf]
    -t <templates>  Specify directory with templates.
                     [default: /var/lib/bithookctl/templates/]
    -p <key>        Specify Bitbucket hook key or an alias for hook key in
                     [hooks] section from configuration file.
                     [default: pre]
    -v              Set template variable <var_name> value to <var_value>.
    -h --help       Show this screen.
    --version       Show version.
```
