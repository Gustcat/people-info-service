package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/Gustcat/people-info-service/internal/lib/filter"
	"github.com/Gustcat/people-info-service/internal/models"
	"github.com/Gustcat/people-info-service/internal/repository"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	tableName = "person"

	idColumn          = "id"
	nameColumn        = "name"
	surnameColumn     = "surname"
	patronymicColumn  = "patronymic"
	genderColumn      = "gender"
	ageColumn         = "age"
	nationalityColumn = "nationality"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(ctx context.Context, DSN string) (*Repo, error) {
	const op = "repository.postgres.NewRepo"

	db, err := pgxpool.Connect(ctx, DSN)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Repo{db: db}, nil
}

func (r *Repo) Close() {
	r.db.Close()
}

func (r *Repo) Create(ctx context.Context, person *models.EnrichmentPerson) (int64, error) {
	const op = "repository.postgres.NewRepo.Create"

	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, surnameColumn, patronymicColumn, genderColumn, ageColumn, nationalityColumn).
		Values(person.Name, person.Surname, person.Patronymic, person.Gender, person.Age, person.Nationality).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("%s: building SQL failed: %w", op, err)
	}

	var id int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, repository.ErrPersonExists
		}

		return 0, fmt.Errorf("%s: executing query failed: %w", op, err)
	}

	return id, nil
}

func (r *Repo) GetByID(ctx context.Context, id int64) (*models.FullPerson, error) {
	const op = "repository.postgres.NewRepo.GetByID"

	builder := sq.Select(idColumn, nameColumn, surnameColumn, patronymicColumn, genderColumn, ageColumn, nationalityColumn).
		From(tableName).
		Where(sq.Eq{idColumn: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: building SQL failed: %w", op, err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}
	if !rows.Next() {
		return nil, repository.ErrPersonNotFound
	}
	defer rows.Close()

	var person models.FullPerson
	err = pgxscan.ScanOne(&person, rows)
	if err != nil {
		return nil, fmt.Errorf("%s: scanning failed: %w", op, err)
	}

	return &person, nil
}

func (r *Repo) List(ctx context.Context, filter *filter.PersonFilter) ([]*models.FullPerson, uint64, error) {
	persons := make([]*models.FullPerson, 0)

	builder := sq.Select("COUNT(*)").
		From("person").
		PlaceholderFormat(sq.Dollar)

	builder = applyPersonFilters(builder, filter)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, 0, err
	}

	var total uint64
	err = r.db.QueryRow(ctx, query, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	if filter.Offset != nil && *filter.Offset >= total {
		return persons, 0, nil
	}

	builder = sq.Select(
		idColumn,
		nameColumn,
		surnameColumn,
		patronymicColumn,
		genderColumn,
		ageColumn,
		nationalityColumn).
		From("person").
		PlaceholderFormat(sq.Dollar).
		OrderBy("id")

	builder = applyPersonFilters(builder, filter)

	if filter.Limit != nil {
		builder = builder.Limit(*filter.Limit)
	}

	if filter.Offset != nil {
		builder = builder.Offset(*filter.Offset)
	}

	query, args, err = builder.ToSql()
	if err != nil {
		return nil, 0, err
	}

	err = pgxscan.Select(ctx, r.db, &persons, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return persons, total, nil
}

func applyPersonFilters(builder sq.SelectBuilder, filter *filter.PersonFilter) sq.SelectBuilder {
	if filter.Name != nil {
		builder = builder.Where(sq.Eq{"name": *filter.Name})
	}
	if filter.Surname != nil {
		builder = builder.Where(sq.Eq{"surname": *filter.Surname})
	}
	if filter.Patronymic != nil {
		builder = builder.Where(sq.Eq{"patronymic": *filter.Patronymic})
	}
	if filter.AgeMin != nil {
		builder = builder.Where(sq.GtOrEq{"age": *filter.AgeMin})
	}
	if filter.AgeMax != nil {
		builder = builder.Where(sq.LtOrEq{"age": *filter.AgeMax})
	}
	if filter.Gender != nil {
		builder = builder.Where(sq.Eq{"surname": *filter.Gender})
	}
	if filter.Nationality != nil {
		builder = builder.Where(sq.Eq{"nationality": *filter.Nationality})
	}

	return builder
}

func (r *Repo) Update(ctx context.Context, id int64, personUpdate *models.PersonUpdate) (*models.FullPerson, error) {
	const op = "repository.postgres.NewRepo.Update"

	builder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	// TODO: разобраться с отсутствием Set
	if personUpdate.Name != nil {
		builder = builder.Set("name", *personUpdate.Name)
	}

	if personUpdate.Surname != nil {
		builder = builder.Set("surname", *personUpdate.Surname)
	}

	if personUpdate.Patronymic != nil {
		builder = builder.Set("patronymic", *personUpdate.Patronymic)
	}

	if personUpdate.Gender != nil {
		builder = builder.Set("gender", *personUpdate.Gender)
	}

	if personUpdate.Age != nil {
		builder = builder.Set("age", *personUpdate.Age)
	}

	if personUpdate.Nationality != nil {
		builder = builder.Set("nationality", *personUpdate.Nationality)
	}

	builder = builder.Suffix("RETURNING id, name, surname, patronymic, gender, age, nationality")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: building SQL failed: %w", op, err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}
	defer rows.Close()

	var person models.FullPerson
	err = pgxscan.ScanOne(&person, rows)
	if err != nil {
		return nil, fmt.Errorf("%s: scanning failed: %w", op, err)
	}

	return &person, nil

}

func (r *Repo) Delete(ctx context.Context, id int64) error {
	const op = "repository.postgres.NewRepo.Delete"

	q := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf("%s: building SQL failed: %w", op, err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: executing query failed: %w", op, err)
	}

	return nil
}
