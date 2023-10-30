package repository_test

//var testPerson = &entity.Person{
//	ID:        1,
//	Email:     "test@test.ru",
//	Phone:     "8999",
//	FirstName: "test",
//}
//func TestPersonRepository_GetAll(t *testing.T) {
//	ListPerson := make([]*entity.Person, 0)
//	ListPerson = append(ListPerson, testPerson)
//	db, mock, err := sqlmock.New()
//	require.NoError(t, err)
//	defer func() {
//		mock.ExpectClose()
//		db.Close()
//	}()
//
//	query := `SELECT (id, email, phone, first_name) FROM persons`
//
//	rows := sqlmock.NewRows([]string{"id", "email", "phone", "first_name"}).
//		AddRow(testPerson.ID, testPerson.Email, testPerson.Phone, testPerson.FirstName)
//
//	mock.ExpectQuery(query).WillReturnRows(rows)
//
//	personRepository := repository.NewPersonRepository()
//	result, err := personRepository.GetPersons(getDb(db))
//
//	assert.NoError(t, err)
//	assert.NotNil(t, result)
//	assert.Equal(t, ListPerson, result)
//
//	require.NoError(t, mock.ExpectationsWereMet())
//}
//
//func TestPersonRepository_GetByID(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	require.NoError(t, err)
//	defer func() {
//		mock.ExpectClose()
//		db.Close()
//	}()
//
//	query := fmt.Sprintf(
//		`SELECT id, email, phone, first_name
//		FROM persons WHERE (id = %v)`,
//		testPerson.ID)
//
//	rows := sqlmock.NewRows([]string{"id", "email", "phone", "first_name"}).
//		AddRow(testPerson.ID, testPerson.Email, testPerson.Phone, testPerson.FirstName)
//
//	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)
//
//	personRepository := repository.NewPersonRepository()
//	result, err := personRepository.GetOnePerson(getDb(db), testPerson.ID)
//
//	assert.NoError(t, err)
//	assert.NotNil(t, result)
//	assert.Equal(t, testPerson, result)
//
//	require.NoError(t, mock.ExpectationsWereMet())
//}
//
//func TestPersonRepository_GetByEmail(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	require.NoError(t, err)
//	defer func() {
//		mock.ExpectClose()
//		db.Close()
//	}()
//
//	query := fmt.Sprintf(
//		`SELECT id, email, phone, first_name
//		FROM persons WHERE (email = '%v')`,
//		testPerson.Email)
//
//	rows := sqlmock.NewRows([]string{"id", "email", "phone", "first_name"}).
//		AddRow(testPerson.ID, testPerson.Email, testPerson.Phone, testPerson.FirstName)
//
//	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)
//
//	personRepository := repository.NewPersonRepository()
//	result, err := personRepository.GetByEmail(getDb(db), testPerson.Email)
//
//	assert.NoError(t, err)
//	assert.NotNil(t, result)
//	assert.Equal(t, testPerson, result)
//
//	require.NoError(t, mock.ExpectationsWereMet())
//}
//
//func TestPersonRepository_Create(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	require.NoError(t, err)
//	defer func() {
//		mock.ExpectClose()
//		db.Close()
//	}()
//
//	query := fmt.Sprintf(
//		`INSERT INTO "persons" ("email","phone","first_name","created_at")
//		VALUES ('%v','%v','%v','%v') RETURNING "id"`,
//		testPerson.Email,
//		testPerson.Phone,
//		testPerson.FirstName,
//		time.Now().Format(time.DateTime))
//
//	mock.ExpectQuery(regexp.QuoteMeta(query)).
//		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
//	personRepository := repository.NewPersonRepository()
//	result, err := personRepository.Create(getDb(db), testPerson)
//
//	assert.NoError(t, err)
//	assert.NotNil(t, result)
//	assert.Equal(t, testPerson, result)
//
//	require.NoError(t, mock.ExpectationsWereMet())
//}
//
//func TestPersonRepository_Update(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	require.NoError(t, err)
//	defer func() {
//		mock.ExpectClose()
//		db.Close()
//	}()
//
//	query := fmt.Sprintf(
//		`UPDATE "persons"
//		SET "email" = '%v', "phone" = '%v', "first_name" = '%v', "updated_at" = '%v'
//       WHERE (id = %v) RETURNING "id"`,
//		testPerson.Email,
//		testPerson.Phone,
//		testPerson.FirstName,
//		time.Now().Format(time.DateTime),
//		testPerson.ID)
//
//	mock.ExpectQuery(regexp.QuoteMeta(query)).
//		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
//	personRepository := repository.NewPersonRepository()
//	result, err := personRepository.Update(getDb(db), testPerson.ID, testPerson)
//
//	assert.NoError(t, err)
//	assert.NotNil(t, result)
//	assert.Equal(t, testPerson, result)
//
//	require.NoError(t, mock.ExpectationsWereMet())
//}
//
//func TestPersonRepository_Delete(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	require.NoError(t, err)
//	defer func() {
//		mock.ExpectClose()
//		db.Close()
//	}()
//
//	query := fmt.Sprintf(`DELETE FROM "persons" WHERE (id = %v)`, testPerson.ID)
//	mock.ExpectExec(regexp.QuoteMeta(query)).WillReturnResult(sqlmock.NewResult(1, 1))
//
//	personRepository := repository.NewPersonRepository()
//	err = personRepository.Delete(getDb(db), testPerson.ID)
//
//	assert.NoError(t, err)
//
//	require.NoError(t, mock.ExpectationsWereMet())
//}
//
//func TestPersonRepository_ParseData(t *testing.T) {
//	personRepository := repository.NewPersonRepository()
//	PersonJson, _ := json.Marshal(testPerson)
//	result, err := personRepository.ParseData(PersonJson)
//	assert.NoError(t, err)
//	assert.NotNil(t, result)
//	assert.Equal(t, testPerson, result)
//}
//
//func getDb(db *sql.DB) context.Context {
//	conn := &dbr.Connection{
//		DB:            db,
//		EventReceiver: &dbr.NullEventReceiver{},
//		Dialect:       dialect.PostgreSQL,
//	}
//	sess := conn.NewSession(nil)
//	ctx := context.WithValue(context.Background(), constants.DbConn, sess)
//	return ctx
//}
