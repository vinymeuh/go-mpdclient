// Copyright 2019 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package mpdclient

import (
	"net"
	"testing"
)

func TestConnections(t *testing.T) {
	server, client := net.Pipe()
	go mockMPDServer(t, server)

	if _, err := mpdConnect(client); err != nil {
		t.Fatalf("connection to mockMPDServer must pass")
	}

	failingTests := []string{"localhost:1", "host.notexists:6600"}
	for _, addr := range failingTests {
		if _, err := NewClient(addr); err == nil {
			t.Fatalf("connection to '%s' must fail", addr)
		}
	}
}

func TestCurrentSongCommand(t *testing.T) {
	server, client := net.Pipe()
	go mockMPDServer(t, server)

	mpc, _ := mpdConnect(client)

	song, err := mpc.CurrentSong()
	if err == nil {
		if currentSongTestData.File != song.File {
			t.Fatalf("File fields mismatch")
		}
		if currentSongTestData.ID != song.ID {
			t.Fatalf("ID fields mismatch")
		}
		if currentSongTestData.Name != song.Name {
			t.Fatalf("Name fields mismatch")
		}
		if currentSongTestData.Pos != song.Pos {
			t.Fatalf("Pos fields mismatch")
		}
		if currentSongTestData.Time != song.Time {
			t.Fatalf("Time fields mismatch")
		}
		if currentSongTestData.Title != song.Title {
			t.Fatalf("Title fields mismatch")
		}
	} else {
		t.Fatalf("mpc.CurrentSong() fails with error %#v", err)
	}
}

func TestStatusCommands(t *testing.T) {
	server, client := net.Pipe()
	go mockMPDServer(t, server)

	mpc, _ := mpdConnect(client)

	status, err := mpc.Status()
	if err == nil {
		if statusTestData.Audio != status.Audio {
			t.Fatalf("Audio fields mismatch")
		}
		if statusTestData.State != status.State {
			t.Fatalf("State fields mismatch")
		}
		if statusTestData.Duration != status.Duration {
			t.Fatalf("Duration fields mismatch")
		}
		if statusTestData.Elapsed != status.Elapsed {
			t.Fatalf("Elapsed fields mismatch")
		}
		if statusTestData.Error != status.Error {
			t.Fatalf("Error fields mismatch")
		}
		if statusTestData.Playlist != status.Playlist {
			t.Fatalf("Playlist fields mismatch")
		}
		if statusTestData.State != status.State {
			t.Fatalf("State fields mismatch")
		}
		if statusTestData.Song != status.Song {
			t.Fatalf("Song fields mismatch")
		}
		if statusTestData.Songid != status.Songid {
			t.Fatalf("Songid fields mismatch")
		}
	} else {
		t.Fatalf("mpc.Status() fails with error %#v", err)
	}

	err = mpc.ClearError()
	if err != nil {
		t.Fatalf("mpc.ClearError() fails with error %#v", err)
	}
}

func TestStatsCommand(t *testing.T) {
	server, client := net.Pipe()
	go mockMPDServer(t, server)

	mpc, _ := mpdConnect(client)

	stats, err := mpc.Stats()
	if err == nil {
		if statsTestData.Artists != stats.Artists {
			t.Fatalf("Artists fields mismatch")
		}
		if statsTestData.Albums != stats.Albums {
			t.Fatalf("Albums fields mismatch")
		}
		if statsTestData.DBplaytime != stats.DBplaytime {
			t.Fatalf("DBplaytime fields mismatch")
		}
		if statsTestData.DBupdate != stats.DBupdate {
			t.Fatalf("DBupdate fields mismatch")
		}
		if statsTestData.Playtime != stats.Playtime {
			t.Fatalf("Playtime fields mismatch")
		}
		if statsTestData.Songs != stats.Songs {
			t.Fatalf("Songs fields mismatch")
		}
		if statsTestData.Uptime != stats.Uptime {
			t.Fatalf("Uptime fields mismatch")
		}
	} else {
		t.Fatalf("mpc.Stats() fails with error %#v", err)
	}
}

func TestControllingPlaybackCommands(t *testing.T) {
	server, client := net.Pipe()
	go mockMPDServer(t, server)

	mpc, _ := mpdConnect(client)

	err := mpc.Next()
	if err != nil {
		t.Fatalf("mpc.Next() fails with error %#v", err)
	}

	err = mpc.Pause(false)
	if err != nil {
		t.Fatalf("mpc.Pause(false) fails with error %#v", err)
	}
	err = mpc.Pause(true)
	if err != nil {
		t.Fatalf("mpc.Pause(true) fails with error %#v", err)
	}

	err = mpc.Play(-1)
	if err != nil {
		t.Fatalf("mpc.Play(-1) fails with error %#v", err)
	}
	err = mpc.Play(0)
	if err != nil {
		t.Fatalf("mpc.Play(0) fails with error %#v", err)
	}
	err = mpc.Play(10)
	if err != nil {
		t.Fatalf("mpc.Play(10) fails with error %#v", err)
	}

	err = mpc.PlayID(-1)
	if err != nil {
		t.Fatalf("mpc.PlayID(-1) fails with error %#v", err)
	}
	err = mpc.PlayID(0)
	if err != nil {
		t.Fatalf("mpc.PlayID(0) fails with error %#v", err)
	}
	err = mpc.PlayID(10)
	if err != nil {
		t.Fatalf("mpc.PlayID(10) fails with error %#v", err)
	}

	err = mpc.Previous()
	if err != nil {
		t.Fatalf("mpc.Previous() fails with error %#v", err)
	}

	err = mpc.Stop()
	if err != nil {
		t.Fatalf("mpc.Stop() fails with error %#v", err)
	}
}

func TestStoredPlaylistsCommands(t *testing.T) {
	server, client := net.Pipe()
	go mockMPDServer(t, server)

	mpc, _ := mpdConnect(client)

	err := mpc.Load("GreatestHits")
	if err != nil {
		t.Fatalf("mpc.Load(\"GreatestHits\") fails with error %#v", err)
	}

	err = mpc.Load("nonexistent")
	if err == nil {
		t.Fatalf("mpc.Load(\"nonexistent\") fails with error %#v", err)
	}
}

func TestConnectionSettingsCommands(t *testing.T) {
	server, client := net.Pipe()
	go mockMPDServer(t, server)

	mpc, _ := mpdConnect(client)

	err := mpc.Ping()
	if err != nil {
		t.Fatalf("mpc.Ping() fails with error %#v", err)
	}

	err = mpc.Close()
	if err != nil {
		t.Fatalf("mpc.Close() fails with error %#v", err)
	}
}
