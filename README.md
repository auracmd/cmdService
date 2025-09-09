# cmdService
**cmdService/library:** Executes system commands, receives their output, and returns it to the orchestrator.

The point of this library/service is to call and execute a command and retrieve the output. And ensure to save/have a history of executed commands until a new topic requested.

History will be used for the AI later to diagnose or analyze on the current specific session/question until its done to start a new one. Its going to be included in each prompt is requested by the user.

### Checklist:
- [X] Execute commands and return result.
- [X] Save commands input and output.
- [X] Able to check current running OS.
