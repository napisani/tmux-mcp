package desc

const PaneOutputDescription = "The active tmux panes in the current tmux server. Each pane has the following attributes:\n" + `
Active         1 if active pane
AtBottom       1 if pane is at the bottom of window
AtLeft         1 if pane is at the left of window
AtRight        1 if pane is at the right of window
AtTop          1 if pane is at the top of window
Bg             Pane background color
Bottom         Bottom of pane
CurrentCommand Current command if available
CurrentPath    Current path if available
Dead           1 if pane is dead
DeadSignal     Exit signal of process in dead pane
DeadStatus     Exit status of process in dead pane
DeadTime       Exit time of process in dead pane
Fg             Pane foreground color
Format         1 if format is for a pane
Height         Height of pane
Id             Unique pane ID
InMode         Number of modes pane is in
Index          Index of pane
InputOff       1 if input to pane is disabled
Last           1 if last pane
Left           Left of pane
Marked         1 if this is the marked pane
MarkedSet      1 if a marked pane is set
Mode           Name of pane mode, if any
Path           Path of pane (can be set by application)
Pid            PID of first process in pane
Pipe           1 if pane is being piped
Right          Right of pane
SearchString   Last search string in copy mode
SesssionName    Name of session containing pane
StartCommand   Command pane started with
StartPath      Path pane started with
Synchronized   1 if pane is synchronized
Tabs           Pane tab positions
Title          Title of pane (can be set by application)
Top            Top of pane
Tty            Pseudo terminal of pane
UnseenChanges  1 if there were changes in pane while in mode
Width          Width of pane
WindowIndex    Index of window containing pane
`

const SessionListDescription = "The sessions in the current tmux server. These are the attributes of each session:\n" + `
Activity          Time of session last activity
Alerts            List of window indexes with alerts
Attached          Number of clients session is attached to
AttachedList      List of clients session is attached to
Created           Time session created
Format            1 if format is for a session
Group             Name of session group
GroupAttached     Number of clients sessions in group are attached to
GroupAttachedList List of clients sessions in group are attached to
GroupList         List of sessions in group
GroupManyAttached 1 if multiple clients attached to sessions in group
GroupSize         Size of session group
Grouped           1 if session in a group
Id                Unique session ID
LastAttached      Time session last attached
ManyAttached      1 if multiple clients attached
Marked            1 if this session contains the marked pane
Name              Name of session
Path              Working directory of session
Stack             Window indexes in most recent order
Windows           Number of windows in session
`

const WindowListDescription = "The active tmux windows in the current tmux server. Each window has the following attributes:\n" + `
Active             1 if window active
ActiveClients      Number of clients viewing this window
ActiveClientsList  List of clients viewing this window
ActiveSessions     Number of sessions on which this window is active
ActiveSessionsList List of sessions on which this window is active
Activity           Time of window last activity
ActivityFlag       1 if window has activity
BellFlag           1 if window has bell
Bigger             1 if window is larger than client
CellHeight         Height of each cell in pixels
CellWidth          Width of each cell in pixels
EndFlag            1 if window has the highest index
Flags              Window flags with # escaped as ##
Format             1 if format is for a window
Height             Height of window
Id                 Unique window ID
Index              Index of window
LastFlag           1 if window is the last used
Layout             Window layout description, ignoring zoomed window panes
Linked             1 if window is linked across sessions
LinkedSessions     Number of sessions this window is linked to
LinkedSessionsList List of sessions this window is linked to
MarkedFlag         1 if window contains the marked pane
Name               Name of window
OffsetX            X offset into window if larger than client
OffsetY            Y offset into window if larger than client
Panes              Number of panes in window
RawFlags           Window flags with nothing escaped
SilenceFlag        1 if window has silence alert
StackIndex         Index in session most recent stack
StartFlag          1 if window has the lowest index
VisibleLayout      Window layout description, respecting zoomed window panes
Width              Width of window
ZoomedFlag         1 if window is zoomed
`
