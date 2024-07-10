// Copyright 2019 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package mpdclient

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"regexp"
	"testing"
)

func (s CurrentSong) String() string {
	return fmt.Sprintf("file: %s\nId: %d\nName: %s\nPos: %d\nTime: %d\nTitle: %s\n",
		s.File, s.ID, s.Name, s.Pos, s.Time, s.Title)
}

var currentSongTestData = CurrentSong{
	File:  "http://novazz.ice.infomaniak.ch/novazz-128.mp3",
	ID:    1,
	Name:  "novazz",
	Pos:   0,
	Time:  0,
	Title: "",
}

func (s Status) String() string {
	return fmt.Sprintf("audio: %s\nbitrate: %d\nduration: %d\nelapsed: %f\nerror: %s\nplaylist: %d\nstate: %s\nsong: %d\nsongid: %d\n",
		s.Audio, s.Bitrate, s.Duration, s.Elapsed, s.Error, s.Playlist, s.State, s.Song, s.Songid)
}

var statusTestData = Status{
	Playlist: 17,
	State:    "play",
	Song:     0,
	Songid:   2,
	Elapsed:  2507.560,
	Duration: 600,
	Bitrate:  128,
	Audio:    "44100:16:2",
	Error:    "",
}

func (s Stats) String() string {
	return fmt.Sprintf("artists: %d\nalbums: %d\ndb_playtime: %d\ndb_update: %d\nplaytime: %d\nsongs: %d\nuptime: %d\n",
		s.Artists, s.Albums, s.DBplaytime, s.DBupdate, s.Playtime, s.Songs, s.Uptime)
}

var statsTestData = Stats{
	Artists:    100,
	Albums:     110,
	DBplaytime: 380000,
	DBupdate:   1550336499,
	Playtime:   400,
	Songs:      1500,
	Uptime:     500,
}

func mockMPDServer(t *testing.T, server net.Conn) {
	fmt.Fprintln(server, "OK MPD 0.19.0")
	go mockMPDHandleConnection(t, server)
}

func mockMPDHandleConnection(t *testing.T, conn net.Conn) {
	simpleCommands := regexp.MustCompile(`^(clearerror|ping|next|previous|stop)\r\n`)
	pauseCommand := regexp.MustCompile(`^pause\s+(0|1)\r\n`)
	playCommand := regexp.MustCompile(`^(play|playid)\s+(-1|\d+)\r\n`)
	loadCommand := regexp.MustCompile(`^load\s+(\w+)\r\n`)
	statusCommand := regexp.MustCompile(`^status\r\n`)
	statsCommand := regexp.MustCompile(`^stats\r\n`)
	currentSongCommand := regexp.MustCompile(`^currentsong\r\n`)

	buf := make([]byte, 4096)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}

		b := bytes.Trim(buf, "\x00")
		// fmt.Println(string(b))

		switch {
		case currentSongCommand.Match(b):
			fmt.Fprintf(conn, "%sOK\n", currentSongTestData)
		case statusCommand.Match(b):
			fmt.Fprintf(conn, "%sOK\n", statusTestData)
		case statsCommand.Match(b):
			fmt.Fprintf(conn, "%sOK\n", statsTestData)
		case simpleCommands.Match(b), pauseCommand.Match(b), playCommand.Match(b):
			fmt.Fprintln(conn, "OK")
		case loadCommand.Match(b):
			m := loadCommand.FindStringSubmatch(string(b))
			if m[1] == "nonexistent" {
				fmt.Fprintln(conn, "ACK [50@0] {load} No such playlist")
			} else {
				fmt.Fprintln(conn, "OK")
			}
		default:
			fmt.Fprintf(conn, "ACK {%s}", string(b))
		}
	}
	conn.Close()
}
