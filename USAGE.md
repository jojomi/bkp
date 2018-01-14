# Usage

`bkp` internally uses [`restic`](https://github.com/restic/restic) for the backup execution.

## Step 1: Create a config file

There is two types of config files that `bkp` uses, namely `targets` and `jobs`.
A target could be your external hard drive or any other backend `restic` currently supports.
The syntax for both is [YAML](https://en.wikipedia.org/wiki/YAML).

The paths searched for config are:

    /etc/bkp
    ~/.bkp
    .bkp (relative to current working dir)

You can have arbitrary folder layout *inside* the `targets` and `jobs` folders. Use this to group and organize your backup jobs.

### target

First create a target, for example in `/etc/bkp/targets/external-hdd.yml`:

``` yaml
name: exthdd

password: thisismybackuppassword

type: local
path: /run/media/johannes/EXTDISK/backup/restic/
```

Then create a first job that uses this target for backup, located in `/etc/bkp/jobs/root.yml`:

```
name: root
description: Backup of the root partition of my personal laptop

# config
source: /
target: exthdd
args:
- "--exclude-file"
- "/etc/bkp/restic-exclude"
- "--exclude"
- "/run/storage"
```

To check the configuration `bkp` sees, you can call the `jobs` command:

    bkp jobs

which would output this:

```
1 job evaluated
"root" to "external" (defined in /etc/bkp/jobs/root.yml)
```

If you make sure your external hard disk is mounted at the path given in its config, the job should be marked as relevant. Check like this:

    bkp jobs --relevant
    1 job evaluated
    1 job relevant
    "Root" to "external" (defined in /etc/bkp/jobs/root.yml)


## Step 2: Create a Backup

Creating a new backup should be as simple as invoking – and it is!

    bkp

This will offer all relevant backup jobs for you to choose from. For automation purposes you can add the `--all` parameter preventing the question.

Be prepared the backup can possibly take a lot of time in case it is the first run. Subsequent runs will be a lot faster and require almost no space apart of changed files.


## Step 3: Test restore

Call `bkp restore root` to mount the target for restore.