-- name: CreateGame :one
insert into games (id, name, description, url, starting_at, ending_at)
values ($1, $2, $3, $4, $5, $6)
returning *;

-- name: GetGames :many
select * from games;

-- name: GetGamesOrderByStartingDate :many
select * from games order by starting_at desc;

-- name: GetGameById :one
select * from games where id = $1;

-- name: GetGameByStartingDate :many
select * from games where date(starting_at) = $1;

-- name: UpdateGame :exec
update games
set name = $2,
    description = $3,
    url = $4,
    starting_at = $5,
    ending_at = $6,
    updated_at = now()
where id = $1;

-- name: DeleteGame :exec
delete from games where id = $1;