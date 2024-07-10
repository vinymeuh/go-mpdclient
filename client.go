// Copyright 2019 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package mpdclient

import (
	"io"
	"net"
	"net/textproto"
	"strconv"
	"strings"
)

// Client implements the MPD protocol
type Client struct {
	addr    string
	conn    *textproto.Conn
	version string // the version of the protocol spoken, not the real version of the daemon
}

// NewClient connects a client to the MPD server listening on tcp addr
func NewClient(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	mpc, err := mpdConnect(conn)
	if mpc != nil {
		mpc.addr = addr
	}
	return mpc, err
}

func mpdConnect(conn io.ReadWriteCloser) (*Client, error) {
	c := &Client{
		conn: textproto.NewConn(conn),
	}
	line, err := c.conn.ReadLine()
	if err != nil {
		return nil, err
	}
	if line[0:6] != "OK MPD" {
		return nil, textproto.ProtocolError(line)
	}
	c.version = strings.Split(line, " ")[2]
	return c, nil
}

// Request execute a MPD command
func (c *Client) Request(format string, args ...interface{}) ([]string, error) {
	id, err := c.conn.Cmd(format, args...)
	if err != nil {
		return nil, err
	}
	c.conn.StartResponse(id)
	defer c.conn.EndResponse(id)

	var response []string
	for {
		line, err := c.conn.ReadLine()
		if err != nil {
			return nil, err
		}
		if line == "OK" {
			return response, nil
		}
		if line[0:3] == "ACK" {
			return nil, textproto.ProtocolError(line)
		}
		response = append(response, line)
	}
}

/***************************/
/** Querying MPD’s status **/
/***************************/

// ClearError clears the current error message in status
func (c *Client) ClearError() error {
	_, err := c.Request("clearerror")
	return err
}

// CurrentSong contains the song info returned by CurrentSong()
type CurrentSong struct {
	Album  string
	Artist string
	File   string
	ID     int
	Name   string
	Pos    int
	Time   int
	Title  string
}

// CurrentSong displays the song info of the current song
func (c *Client) CurrentSong() (*CurrentSong, error) {
	resp, err := c.Request("currentsong")
	if err != nil {
		return nil, err
	}

	var s CurrentSong
	for _, line := range resp {
		elm := strings.Split(line, ": ")
		if len(elm) < 2 {
			return nil, textproto.ProtocolError(line)
		}
		switch elm[0] {
		case "Album":
			s.Album = elm[1]
		case "Artist":
			s.Artist = elm[1]
		case "file":
			s.File = elm[1]
		case "Id":
			s.ID, _ = strconv.Atoi(elm[1])
		case "Name":
			s.Name = elm[1]
		case "Pos":
			s.Pos, _ = strconv.Atoi(elm[1])
		case "Time":
			s.Time, _ = strconv.Atoi(elm[1])
		case "Title":
			s.Title = elm[1]
		}
	}
	return &s, nil
}

// Idle waits until there is a noteworthy change in one or more of MPD’s subsystems
// While a client is waiting for idle results, the server disables timeouts,
// allowing a client to wait for events as long as mpd runs.
// Change events accumulate, even while the connection is not in “idle” mode, no events gets
// lost while the client is doing something else with the connection.
// If an event had already occurred since the last call, the new idle command will return immediately.
func (c *Client) Idle(subsystems string) ([]string, error) {
	resp, err := c.Request("idle %s", subsystems)
	if err != nil {
		return nil, err
	}

	events := make([]string, 0, 6)
	i := 0
	for _, line := range resp {
		elm := strings.Split(line, ": ")
		if len(elm) < 2 {
			return nil, textproto.ProtocolError(line)
		}
		switch elm[0] {
		case "changed":
			events = append(events, elm[1])
			i++
		}
	}
	return events, nil
}

// NoIdle cancels Idle Command (no other commands are allowed)
func (c *Client) NoIdle() error {
	_, err := c.Request("noidle")
	return err
}

// Status contains current status of the player returned by Status()
// Note: some fields returned by MPD are ignored
type Status struct {
	Audio          string
	Bitrate        int
	Duration       int
	Elapsed        float64
	Error          string
	Playlist       int
	PlaylistLength int
	Song           int
	Songid         int
	State          string
}

// Status reports the current status of the player and the volume level
func (c *Client) Status() (*Status, error) {
	resp, err := c.Request("status")
	if err != nil {
		return nil, err
	}

	var s Status
	for _, line := range resp {
		elm := strings.Split(line, ": ")
		if len(elm) < 2 {
			return nil, textproto.ProtocolError(line)
		}
		switch elm[0] {
		case "audio":
			s.Audio = elm[1]
		case "bitrate":
			s.Bitrate, _ = strconv.Atoi(elm[1])
		case "duration":
			s.Duration, _ = strconv.Atoi(elm[1])
		case "elapsed":
			s.Elapsed, _ = strconv.ParseFloat(elm[1], 64)
		case "error":
			s.Error = elm[1]
		case "playlist":
			s.Playlist, _ = strconv.Atoi(elm[1])
		case "playlistlength":
			s.PlaylistLength, _ = strconv.Atoi(elm[1])
		case "song":
			s.Song, _ = strconv.Atoi(elm[1])
		case "songid":
			s.Songid, _ = strconv.Atoi(elm[1])
		case "state":
			s.State = elm[1]
		}
	}
	return &s, nil
}

// Stats contains server statistics returned by Stats()
type Stats struct {
	Artists    int
	Albums     int
	DBplaytime int
	DBupdate   int
	Playtime   int
	Songs      int
	Uptime     int
}

// Stats displays server statistics
func (c *Client) Stats() (*Stats, error) {
	resp, err := c.Request("stats")
	if err != nil {
		return nil, err
	}

	var s Stats
	for _, line := range resp {
		elm := strings.Split(line, ": ")
		if len(elm) < 2 {
			return nil, textproto.ProtocolError(line)
		}
		switch elm[0] {
		case "artists":
			s.Artists, _ = strconv.Atoi(elm[1])
		case "albums":
			s.Albums, _ = strconv.Atoi(elm[1])
		case "db_playtime":
			s.DBplaytime, _ = strconv.Atoi(elm[1])
		case "db_update":
			s.DBupdate, _ = strconv.Atoi(elm[1])
		case "playtime":
			s.Playtime, _ = strconv.Atoi(elm[1])
		case "songs":
			s.Songs, _ = strconv.Atoi(elm[1])
		case "uptime":
			s.Uptime, _ = strconv.Atoi(elm[1])
		}
	}
	return &s, nil
}

/**************************/
/** Controlling playback **/
/**************************/

// Next plays next song in the playlist
func (c *Client) Next() error {
	_, err := c.Request("next")
	return err
}

// Pause toggles pause/resumes playing
func (c *Client) Pause(pause bool) error {
	var err error
	switch pause {
	case true:
		_, err = c.Request("pause 1")
	case false:
		_, err = c.Request("pause 0")
	}
	return err
}

// Play begins playing the playlist at song number songpos.
// If songpos equal -1, starts playing at the current position
// in the playlist: this is not documented in the protocol but
// successfully tested with MPD 0.19.0
func (c *Client) Play(songpos int) error {
	_, err := c.Request("play %d", songpos)
	return err
}

// PlayID begins playing the playlist at song number songpos.
// If songpos equal -1, starts playing at the current position
// in the playlist: this is not documented in the protocol but
// successfully tested with MPD 0.19.0
func (c *Client) PlayID(songid int) error {
	_, err := c.Request("playid %d", songid)
	return err
}

// Previous plays previous song in the playlist
func (c *Client) Previous() error {
	_, err := c.Request("previous")
	return err
}

// Stop stops playing
func (c *Client) Stop() error {
	_, err := c.Request("stop")
	return err
}

/**********************/
/** Stored playlists **/
/**********************/

// Load loads the playlist into the current queue
// Note: range [start:end] currently not supported
func (c *Client) Load(name string) error {
	_, err := c.Request("load %s", name)
	return err
}

/*************************/
/** Connection settings **/
/*************************/

// Close the connection with MPD
func (c *Client) Close() error {
	if c.conn != nil {
		// clients should not use "close" command; instead, they should just close the socket
		err := c.conn.Close()
		c.conn = nil
		c.version = ""
		if err != nil {
			return err
		}
	}
	return nil
}

// Ping does nothing but return “OK”
func (c *Client) Ping() error {
	_, err := c.Request("ping")
	return err
}
