package db

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

//TODO: Wrap errors

const dbTimeout = time.Second * 10
const (
	InWork      = "В обработке"
	InExecution = "В исполнении"
	Rejected    = "Отменен"
	Uploaded    = "Загружен"
)

type Service interface {
	NewUser(string, string, string) error
	GetUserId(string, string) (int, error)
	GetUser(int) (*User, error)
	GetReleaseByUserId(int) ([]*Release, error)
	GetReleaseById(int) (*Release, error)
	GetTrackByReleaseId(int) ([]*Track, error)
}

type service struct {
	db *sql.DB
}

func New(connStr string) (Service, error) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("open sql connection error: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping db error: %w", err)
	}
	return &service{db: db}, nil
}

func (s *service) NewUser(email, username, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	password = hashPass(password)
	query := fmt.Sprintf("insert into users (email, username, password) values ('%s', '%s', '%s')",
		email, username, password)
	fmt.Println(query)
	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("error user insert: %w", err)
	}

	return nil
}

func (s *service) GetUserId(username, password string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	password = hashPass(password)
	rows, err := s.db.QueryContext(ctx,
		fmt.Sprintf("select user_id from users where username='%s' and password ='%s'", username, password))
	if err != nil {
		return 0, fmt.Errorf("error user select: %w", err)
	}
	var id int
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return 0, fmt.Errorf("error user scan by id: %w", err)
		}
	}
	if err != nil {
		return 0, fmt.Errorf("error getting user last id: %w", err)
	}
	return id, nil
}

func (s *service) GetUser(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	rows, err := s.db.QueryContext(ctx,
		fmt.Sprintf("select user_id, email, username, password from users where user_id='%d'", id))
	if err != nil {
		return nil, fmt.Errorf("error user select by id: %w", err)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("row error: %w", err)
	}
	var user User
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Email, &user.Username, &user.Password)
		if err != nil {
			return nil, fmt.Errorf("eror user scan: %w", err)
		}

	}
	return &user, nil
}

func (s *service) NewRelease(userId int, name, authors, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := fmt.Sprintf("insert into releases (user_id, name, authors, status) values ('%d', '%s', '%s', '%s')",
		userId, name, authors, status)
	fmt.Println(query)
	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("error release insert: %w", err)
	}

	return nil
}

func (s *service) GetReleaseId(name, authors string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	rows, err := s.db.QueryContext(ctx,
		fmt.Sprintf("select release_id from releases where name='%s' and authors ='%s'", name, authors))
	if err != nil {
		return 0, fmt.Errorf("error release select: %w", err)
	}
	var id int
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return 0, fmt.Errorf("error release scan by id: %w", err)
		}
	}
	if err != nil {
		return 0, fmt.Errorf("error getting release last id: %w", err)
	}
	return id, nil
}

func (s *service) GetReleaseByUserId(userId int) ([]*Release, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	rows, err := s.db.QueryContext(ctx,
		fmt.Sprintf("select release_id, cover, name, authors, status, release_date from releases where user_id='%d'",
			userId))
	if err != nil {
		return nil, fmt.Errorf("error releases select by user id: %w", err)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("row error: %w", err)
	}
	var releases []*Release
	for rows.Next() {
		var release Release
		err := rows.Scan(&release.Id, &release.Cover, &release.Name, &release.Authors, &release.Status, &release.ReleaseDate)
		if err != nil {
			return nil, fmt.Errorf("eror release scan: %w", err)
		}
		release.ReleaseDate = cutDate(release.ReleaseDate)
		releases = append(releases, &release)
	}
	return releases, nil
}
func (s *service) GetReleaseById(releaseId int) (*Release, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	rows, err := s.db.QueryContext(ctx,
		fmt.Sprintf("select release_id, cover, name, authors, status, release_date from releases where release_id='%d'",
			releaseId))
	if err != nil {
		return nil, fmt.Errorf("error releases select by user id: %w", err)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("row error: %w", err)
	}
	var release Release
	for rows.Next() {
		err := rows.Scan(&release.Id, &release.Cover, &release.Name, &release.Authors, &release.Status, &release.ReleaseDate)
		if err != nil {
			return nil, fmt.Errorf("eror release scan: %w", err)
		}
		release.ReleaseDate = cutDate(release.ReleaseDate)
	}
	return &release, nil
}

func (s *service) NewTrack(releaseId int, name, authors string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := fmt.Sprintf("insert into tracks (release_id, name, authors) values ('%d', '%s', '%s')",
		releaseId, name, authors)
	fmt.Println(query)
	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("error track insert: %w", err)
	}

	return nil
}

func (s *service) GetTrackId(name, authors string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	rows, err := s.db.QueryContext(ctx,
		fmt.Sprintf("select track_id from tracks where name='%s' and authors ='%s'", name, authors))
	if err != nil {
		return 0, fmt.Errorf("error track select: %w", err)
	}
	var id int
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return 0, fmt.Errorf("error track scan by id: %w", err)
		}
	}
	if err != nil {
		return 0, fmt.Errorf("error getting track last id: %w", err)
	}
	return id, nil
}

func (s *service) GetTrackByReleaseId(releaseId int) ([]*Track, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	rows, err := s.db.QueryContext(ctx,
		fmt.Sprintf("select track_id, name, authors from tracks where release_id='%d'", releaseId))
	if err != nil {
		return nil, fmt.Errorf("error tracks select by release id: %w", err)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("row error: %w", err)
	}
	var tracks []*Track
	for rows.Next() {
		var track Track
		err := rows.Scan(&track.Id, &track.Name, &track.Authors)
		if err != nil {
			return nil, fmt.Errorf("eror track scan: %w", err)
		}
		tracks = append(tracks, &track)
	}
	return tracks, nil
}

func hashPass(pass string) string {
	hashPass := md5.New()
	hashPass.Write([]byte(pass))
	return base64.URLEncoding.EncodeToString(hashPass.Sum(nil))
}

func cutDate(date string) string {
	idx := strings.Index(date, "T")
	if idx >= len(date) {
		return date
	}
	return date[:idx]
}
