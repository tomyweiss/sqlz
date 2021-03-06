package sqlz

import "testing"

func TestInsert(t *testing.T) {
	runTests(t, func(dbz *DB) []test {
		return []test{
			{
				"simple insert",
				dbz.InsertInto("table").Columns("id", "name", "date").Values(1, "My Name", 96969696),
				"INSERT INTO table (id, name, date) VALUES (?, ?, ?)",
				[]interface{}{1, "My Name", 96969696},
			},

			{
				"insert with value map",
				dbz.InsertInto("table").ValueMap(map[string]interface{}{"id": 1, "name": "My Name"}),
				"INSERT INTO table (id, name) VALUES (?, ?)",
				[]interface{}{1, "My Name"},
			},

			{
				"insert with returning clause",
				dbz.InsertInto("table").Columns("one", "two").Values(1, 2).Returning("id"),
				"INSERT INTO table (one, two) VALUES (?, ?) RETURNING id",
				[]interface{}{1, 2},
			},

			{
				"insert with on conflict do nothing clause",
				dbz.InsertInto("table").Columns("one", "two").Values(1, 2).OnConflictDoNothing(),
				"INSERT INTO table (one, two) VALUES (?, ?) ON CONFLICT DO NOTHING",
				[]interface{}{1, 2},
			},

			{
				"insert rows from a select query",
				dbz.InsertInto("table").Columns("one", "two").FromSelect(
					dbz.Select("*").From("table2"),
				),
				"INSERT INTO table (one, two) SELECT * FROM table2",
				[]interface{}{},
			},

			{
				"insert with on conflict do update",
				dbz.InsertInto("table").Columns("name").Values("My Name").
					OnConflict(
						OnConflict("name", "something_else").
							DoUpdate().
							Set("update_date", 55151515).
							SetMap(map[string]interface{}{
								"name":    "My Name Again",
								"address": "Some Address",
							})),
				"INSERT INTO table (name) VALUES (?) ON CONFLICT (name, something_else) DO UPDATE SET update_date = ?, address = ?, name = ?",
				[]interface{}{"My Name", 55151515, "Some Address", "My Name Again"},
			},

			{
				"insert with on conflict do update conditional set",
				dbz.InsertInto("table").Columns("name").Values("My Name").
					OnConflict(
						OnConflict("name").
							DoUpdate().
							Set("update_date", 55151515).
							SetIf("true", 1, true).
							SetIf("false", 0, false),
					),
				"INSERT INTO table (name) VALUES (?) ON CONFLICT (name) DO UPDATE SET update_date = ?, true = ?",
				[]interface{}{"My Name", 55151515, 1},
			},

			{
				"insert or ignore",
				dbz.InsertInto("table").OrIgnore().Columns("id", "name", "date").Values(1, "My Name", 96969696),
				"INSERT OR IGNORE INTO table (id, name, date) VALUES (?, ?, ?)",
				[]interface{}{1, "My Name", 96969696},
			},
			{
				"insert or replace",
				dbz.InsertInto("table").OrReplace().Columns("id", "name", "date").Values(1, "My Name", 96969696),
				"INSERT OR REPLACE INTO table (id, name, date) VALUES (?, ?, ?)",
				[]interface{}{1, "My Name", 96969696},
			},

			{
				"insert with multiple values",
				dbz.InsertInto("table").Columns("id", "name").
					ValueMultiple([][]interface{}{{1, "My Name"}, {2, "John"}, {3, "Golang"}}),
				"INSERT INTO table (id, name) VALUES (?, ?), (?, ?), (?, ?)",
				[]interface{}{1, "My Name", 2, "John", 3, "Golang"},
			},
		}
	})
}
