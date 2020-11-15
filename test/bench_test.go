package test

// func BenchmarkCreate(b *testing.B) {
// 	db, err := initSqliteDB()
// 	if err != nil {
// 		b.Fatalf("err:%+v", err)
// 	}
// 	factory.Opt().SetDB(db, "sqlite3")
// 	locationFactory := factory.New(
// 		func() interface{} { return &Location{} },
// 		attr.Seq("ID", 1, "id"),
// 		attr.Str("Address", randomdata.Address, "address"),
// 	).Table("locations")

// 	homeFactory := factory.New(
// 		func() interface{} { return &Home{} },
// 		attr.Seq("ID", 1, "id"),
// 	).Columns(factory.Col("HostID", "host_id")).FAssociate("Location", locationFactory, 1, true, nil, "location_id").Table("homes")

// 	userFactory := factory.New(
// 		func() interface{} { return &User{CreatedAt: time.Now()} },
// 		attr.Seq("ID", 1, "id"),
// 		attr.Str("Username", genutil.NameRander(3), "username"),
// 		attr.Int("Age", genutil.IntRander(25, 50), "age"),
// 	).FAssociate("Home", homeFactory, 1, false, func(data, depend interface{}) error {
// 		user := data.(*User)
// 		home := depend.(*Home)
// 		home.HostID = user.ID
// 		return nil
// 	}).Table("users")

// 	for i := 0; i < b.N; i++ {
// 		_, err := userFactory.Insert()
// 		if err != nil {
// 			b.Fatalf("err:%+v", err)
// 		}
// 	}
// }
