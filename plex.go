package main

import (
	"fmt"
	"github.com/jyggen/go-plex-client"
	"strings"
	"time"
)

type Plex struct {
	connection *plex.Plex
}

func NewPlex(baseUrl string, token string) (*Plex, error) {
	connection, err := plex.New(strings.TrimRight(baseUrl, "/"), token)

	if err != nil {
		return &Plex{}, err
	}

	connection.HTTPClient.Timeout = time.Second * 10

	_, err = connection.Test()

	if err != nil {
		return &Plex{}, err
	}

	return &Plex{
		connection: connection,
	}, nil
}

func (p *Plex) Shows() ([]*Show, error) {
	shows := make([]*Show, 0)
	libraries, err := p.connection.GetLibraries()

	if err != nil {
		return shows, fmt.Errorf("unable to retrieve libraries: %w", err)
	}

	for _, library := range libraries.MediaContainer.Directory {
		if library.Type != "show" {
			continue
		}

		content, err := p.connection.GetLibraryContent(library.Key, "")

		if err != nil {
			return shows, fmt.Errorf("unable to retrieve content for library %v: %w", library.Key, err)
		}

		for _, show := range content.MediaContainer.Metadata {
			shows = append(shows, &Show{
				connection: p.connection,
				metadata:   show,
			})
		}
	}

	return shows, nil
}

type Show struct {
	connection *plex.Plex
	metadata   plex.Metadata
}

func (s *Show) Seasons() ([]*Season, error) {
	seasons := make([]*Season, 0)
	children, err := s.connection.GetMetadataChildren(s.metadata.RatingKey)

	if err != nil {
		return seasons, err
	}

	for _, season := range children.MediaContainer.Metadata {
		seasons = append(seasons, &Season{
			connection: s.connection,
			metadata:   season,
		})
	}

	return seasons, nil
}

func (s *Show) RatingKey() string {
	return s.metadata.RatingKey
}

func (s *Show) SortTitle() string {
	if s.metadata.TitleSort == "" {
		return s.metadata.Title
	}

	return s.metadata.TitleSort
}

func (s *Show) Title() string {
	return s.metadata.Title
}

type Season struct {
	connection *plex.Plex
	metadata   plex.Metadata
}

func (s *Season) Episodes() ([]*Episode, error) {
	episodes := make([]*Episode, 0)
	children, err := s.connection.GetMetadataChildren(s.metadata.RatingKey)

	if err != nil {
		return episodes, err
	}

	for _, episode := range children.MediaContainer.Metadata {
		metadata, err := s.connection.GetMetadata(episode.RatingKey)

		if err != nil {
			return episodes, err
		}

		episodes = append(episodes, &Episode{
			connection: s.connection,
			metadata:   metadata.MediaContainer.Metadata[0],
		})
	}

	return episodes, nil
}

func (e *Season) Number() int {
	return int(e.metadata.Index)
}

type Episode struct {
	connection *plex.Plex
	metadata   plex.Metadata
}

func (e *Episode) HasIntroMarker() bool {
	for _, marker := range e.metadata.Marker {
		if marker.Type == "intro" {
			return true
		}
	}

	return false
}

func (e *Episode) Number() int {
	return int(e.metadata.Index)
}
