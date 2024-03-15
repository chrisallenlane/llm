llm
===
**This is alpha software. It is under active development, and its interface is
unstable.**


Overview
--------
`llm` is a CLI wrapper for the OpenAI API. (Support for locally-hosted LLMs
with a OpenAI-conformant interface is coming soon). Its purpose is to make LLM
integrations convenient to those with a shell-centric workflow.

`llm` allows for the creation of multiple long-lived, user-created "sessions"
(which correspond with a "chat" in the GPT-4 web interface), each of which may
have a distinct "system" prompt.

Conversation history is stored in a local SQLite3 database, and thus searchable
locally. Additionally, syntax highlighting is enabled by default.


Features
--------
1. Shell-centric workflow that integrates well with shell-based text editors like `vim` and `emacs`
2. Ability to store and search conversation history locally
3. Syntax highlighting
4. Ability to "rewind" and "copy" sessions
5. Ability to "stack" and batch messages


Usage
-----
```sh
Usage:
    llm [options] db (path|destroy)
    llm [options] (q|question) (add|ask) [<msg>]
    llm [options] (q|question) (edit|rm) <id>
    llm [options] (q|question) send
    llm [options] (s|session) (new|edit|view|rm|use) <name>
    llm [options] (s|session) cp <orig> <name> [--all]
    llm [options] (s|session) log [<num>]
    llm [options] (s|session) ls
    llm [options] (s|session) rewind <id>
    llm [options] (s|session) search <string>
    llm [options] (s|session) truncate
```

### Example workflows ###
Simple questions may be asked directly inline:

```sh
$ llm q ask 'Please demonstrate how to install a package on Debian Linux.'
```

Longer questions may be composed in your preferred editor:
```sh
$ llm q ask
# This will open $EDITOR or $VISUAL. Compose your question in the shell, then
save the document. On write and close, the question will be sent to the OpenAI
API.
```
