package models

import (
	"context"
	"database/sql"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

// 指定idのmovieかerrorを返すメソッド(DBModelのポインタレシーバ)
func (m *DBModel) Get(id int) (*Movie, error) {
	// 3sでタイムアウトする
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 指定したIDのmoviesを取得するクエリ
	query := `select id, title, description, year, release_date, runtime, rating, mpaa_rating,
							created_at, updated_at from movies where id = $1
	`

	// 指定したidのmoviesを取得する(1行)
	row := m.DB.QueryRowContext(ctx, query, id)

	var movie Movie

	// クエリの結果をmovieに割り当てる
	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.Description,
		&movie.Year,
		&movie.ReleaseDate,
		&movie.Runtime,
		&movie.Rating,
		&movie.MPAARating,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// 指定したmovie_idのgenresを取得するクエリ
	query = `select
						mg.id, mg.movie_id, mg.genre_id, g.genre_name
					from
						movies_genres mg
						left join genres g on (g.id = mg.genre_id)
					where
						mg.movie_id = $1
	`

	// 指定したmovie_idのgenresを取得する(複数行)
	rows, _ := m.DB.QueryContext(ctx, query, id)
	defer rows.Close()

	genres := make(map[int]string)
	for rows.Next() {
		var mg MovieGenre
		err := rows.Scan(
			&mg.ID,
			&mg.MovieID,
			&mg.GenreID,
			&mg.Genre.GenreName,
		)
		if err != nil {
			return nil, err
		}
		genres[mg.ID] = mg.Genre.GenreName
	}

	movie.MovieGenre = genres

	return &movie, nil
}

// すべてのmovieかerrorを返すメソッド(DBModelのポインタレシーバ)
func (m *DBModel) All() ([]*Movie, error) {
	// 3sでタイムアウトする
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, title, description, year, release_date, runtime, rating, mpaa_rating,
							created_at, updated_at from movies order by title
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []*Movie

	for rows.Next() {
		var movie Movie
		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Description,
			&movie.Year,
			&movie.ReleaseDate,
			&movie.Rating,
			&movie.Runtime,
			&movie.MPAARating,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// get genres, if any
		genreQuery := `select
			mg.id, mg.movie_id, mg.genre_id, g.genre_name
		from
			movies_genres mg
			left join genres g on (g.id = mg.genre_id)
		where
			mg.movie_id = $1
		`

		genreRows, _ := m.DB.QueryContext(ctx, genreQuery, movie.ID)
	
		genres := make(map[int]string)
		for genreRows.Next() {
			var mg MovieGenre
			err := genreRows.Scan(
				&mg.ID,
				&mg.MovieID,
				&mg.GenreID,
				&mg.Genre.GenreName,
			)
			if err != nil {
				return nil, err
			}
			genres[mg.ID] = mg.Genre.GenreName
		}
		genreRows.Close()

		movie.MovieGenre = genres
		movies = append(movies, &movie)
	}


	return movies, nil
}