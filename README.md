# gogo-factory


## Install


## Getting Started
- [Defining factories](#defining-factories)
  - [Factory name and attributes](#factory-name-and-attributes)
  - [Customize attribute with other fields](#customize-attribute-with-other-fields)
- [Building objects](#building-objects)
  - [Building strategies](#building-strategies)
  - [Customize insert function](#customize-insert-function)
- [Inheritance and Association](#inheritance-and-association)
  - [Inherit factory](#inherit-factory)
  - [Factories association](#factories-association)


#### Defining factories

##### Factory name and attributes

```go
type Gender int8

const (
	Male Gender = iota + 1
	Female
)

type Phone string

type User struct {
	ID        int64
	Name      string
	Gender    Gender
	Phone     Phone
	Address   sql.NullString
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}


```

```go
userFactory := factory.New(
	func() interface{} { return &User{} },
	attr.Seq("ID", 1),
	attr.Str("Name", randutil.NameRander(3)),
	attr.Int("Gender", randutil.IntRander(1, 2)),
	attr.Str("Phone", randomdata.PhoneNumber),
	attr.Str("Address", randomdata.Address),
	attr.Time("CreatedAt", randutil.NowRander()),
	attr.Time("UpdatedAt", randutil.TimeRander(time.Now(), time.Now().Add(30*time.Hour))),
)

for i := 0; i < 5; i++ {
  user := userFactory.MustBuild().(*User)
  spew.Printf("user_%d: %+v\n", user.ID, user)
}
```

output:
```shell
{ID:1 Name:Sophia, Jackson Gender:2 Phone:+685 75 86 973 6 1 20 Address:{String:78 Lincoln Circle,
San Martin, NM, 74682 Valid:true} CreatedAt:2020-11-03 14:29:52.91313 +0800 CST m=+0.001783291 UpdatedAt:{Time:2020-11-03 14:29:52 +0800 CST Valid:true}}
{ID:2 Name:Elizabeth, Anderson Gender:1 Phone:+44 627 140670 117 Address:{String:62 Madison Ter,
New Deal, CO, 23106 Valid:true} CreatedAt:2020-11-03 14:29:52.913344 +0800 CST m=+0.001997691 UpdatedAt:{Time:2020-11-03 14:29:52 +0800 CST Valid:true}}
{ID:3 Name:Jacob, Taylor Gender:1 Phone:+850 1 676199016 5 Address:{String:71 Lincoln Circle,
Bury, WV, 31262 Valid:true} CreatedAt:2020-11-03 14:29:52.913389 +0800 CST m=+0.002042964 UpdatedAt:{Time:2020-11-03 14:29:52 +0800 CST Valid:true}}
{ID:4 Name:Aubrey, Miller Gender:1 Phone:+1 684 2 6 8251 0 1 42 Address:{String:9 Madison St,
Derby Center, NH, 06268 Valid:true} CreatedAt:2020-11-03 14:29:52.913429 +0800 CST m=+0.002082182 UpdatedAt:{Time:2020-11-03 14:29:52 +0800 CST Valid:true}}
{ID:5 Name:Ella, Smith Gender:1 Phone:+234 82 40188507 1 Address:{String:92 Franklin Ave,
Northleach, KY, 19480 Valid:true} CreatedAt:2020-11-03 14:29:52.913464 +0800 CST m=+0.002117516 UpdatedAt:{Time:2020-11-03 14:29:52 +0800 CST Valid:true}}
```

##### Customize attribute with other fields

```go
userFactory := factory.New(
  func() interface{} { return &User{} },
  attr.Seq("ID", 1),
  attr.Str("Name", randutil.NameRander(3)).Process(func(a attr.Attributer) error {
    user := a.GetObject().(*User)
    name := a.GetVal().(string)
    name = "username-" + strconv.Itoa(int(user.ID))
    return a.SetVal(name)
  }),
  attr.Int("Gender", randutil.IntRander(1, 2)),
  attr.Attr("Phone", func() interface{} { return Phone(randomdata.PhoneNumber()) }),
  attr.Str("Address", randomdata.Address),
  attr.Time("CreatedAt", randutil.NowRander()),
  attr.Time("UpdatedAt", randutil.TimeRander(time.Now(), time.Now().Add(30*time.Hour))),
)

for i := 0; i < 5; i++ {
  user := userFactory.MustBuild().(*User)
  spew.Printf("user_%d: %+v\n", user.ID, user)
}
```

output:
```shell
{ID:1 Name:username-1 Gender:2 Phone:+1 06 543 1299 1249 Address:{String:76 Franklin St,
Berkhamsted, MI, 82996 Valid:true} CreatedAt:2020-11-03 14:29:52.913504 +0800 CST m=+0.002157066 UpdatedAt:{Time:2020-11-03 14:29:52 +0800 CST Valid:true}}
{ID:2 Name:username-2 Gender:1 Phone:+423 7 1445692 44 6 Address:{String:62 Adams Pkwy,
Derby Center, SC, 33316 Valid:true} CreatedAt:2020-11-03 14:29:52.913539 +0800 CST m=+0.002192055 UpdatedAt:{Time:2020-11-03 14:29:52 +0800 CST Valid:true}}
{ID:3 Name:username-3 Gender:1 Phone:+965 34 4948425 30 Address:{String:9 Franklin Rd,
San Martin, NC, 28261 Valid:true} CreatedAt:2020-11-03 14:29:52.913573 +0800 CST m=+0.002226672 UpdatedAt:{Time:2020-11-03 14:29:52 +0800 CST Valid:true}}
{ID:4 Name:username-4 Gender:1 Phone:+675 1 3350381 334 Address:{String:33 Lincoln Blvd,
Ransom Canyon, UT, 55085 Valid:true} CreatedAt:2020-11-03 14:29:52.91361 +0800 CST m=+0.002263674 UpdatedAt:{Time:2020-11-03 14:29:52 +0800 CST Valid:true}}
{ID:5 Name:username-5 Gender:2 Phone:+974 1 625192933 7 Address:{String:64 Washington Ct,
Brandwell, NE, 95480 Valid:true} CreatedAt:2020-11-03 14:29:52.913652 +0800 CST m=+0.002305552 UpdatedAt:{Time:2020-11-03 14:29:52 +0800 CST Valid:true}}
```

#### Building objects

##### Building strategies

build and insert

```sql
CREATE TABLE IF NOT EXISTS `users` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `name` VARCHAR(64) NULL,
    `address` INTEGER NULL,
    `gender` INTEGER NULL,
    `phone` VARCHAR(30) NULL,
    `created_at` DATE NULL,
    `updated_at` DATE NULL
);
```

```go
// setup global sql.DB instance
db, err := initSqliteDB()
if err != nil {
  log.Panicf("err:%+v", err)
}
factory.DB(db, "sqlite3")

userFactory := factory.New(
  func() interface{} { return &User{} },
  attr.Seq("ID", 1, "id"),
  attr.Str("Name", randutil.NameRander(3), "name"),
  attr.Int("Gender", randutil.IntRander(1, 2), "gender"),
  attr.Str("Phone", randomdata.PhoneNumber, "phone"),
  attr.Str("Address", randomdata.Address, "address"),
  attr.Time("CreatedAt", randutil.NowRander(), "created_at"),
  attr.Time("UpdatedAt", randutil.TimeRander(time.Now(), time.Now().Add(30*time.Hour)), "updated_at"),
).Table("users")

for i := 0; i < 5; i++ {
  user := userFactory.MustInsert().(*User)
  spew.Printf("user_%d: %+v\n", user.ID, user)
}
```

omit fields

```go
userFactory := factory.New(
  func() interface{} { return &User{} },
  attr.Seq("ID", 1, "id"),
  attr.Str("Name", randutil.NameRander(3), "name"),
  attr.Int("Gender", randutil.IntRander(1, 2), "gender"),
  attr.Str("Phone", randomdata.PhoneNumber, "phone"),
  attr.Str("Address", randomdata.Address, "address"),
  attr.Time("CreatedAt", randutil.NowRander(), "created_at"),
  attr.Time("UpdatedAt", randutil.TimeRander(time.Now(), time.Now().Add(30*time.Hour)), "updated_at"),
).Table("users")

for i := 0; i < 5; i++ {
  user := userFactory.Omit("Address").MustInsert().(*User) 
  spew.Printf("%+v\n", user)
}

user := userFactory.MustInsert().(*User) 
spew.Printf("%+v\n", user)
```

output

```shell
# address is empty
{ID:1 Name:Addison, Wilson Gender:2 Phone:+60 83 398 7573466 Address:{String: Valid:false} CreatedAt:2020-11-03 16:41:23.64853 +0800 CST m=+0.006121730 UpdatedAt:{Time:2020-11-03 16:41:23 +0800 CST Valid:true}}
{ID:2 Name:Mason, Thompson Gender:2 Phone:+375 12 4381 995 3 0 Address:{String: Valid:false} CreatedAt:2020-11-03 16:41:23.64942 +0800 CST m=+0.007011692 UpdatedAt:{Time:2020-11-03 16:41:23 +0800 CST Valid:true}}
{ID:3 Name:Lily, Wilson Gender:1 Phone:+421 442 4930731 7 Address:{String: Valid:false} CreatedAt:2020-11-03 16:41:23.650046 +0800 CST m=+0.007638218 UpdatedAt:{Time:2020-11-03 16:41:23 +0800 CST Valid:true}}
{ID:4 Name:Emily, Robinson Gender:2 Phone:+7 77 208172 1 1 61 0 Address:{String: Valid:false} CreatedAt:2020-11-03 16:41:23.650647 +0800 CST m=+0.008238494 UpdatedAt:{Time:2020-11-03 16:41:23 +0800 CST Valid:true}}
{ID:5 Name:Zoey, Garcia Gender:2 Phone:+244 0 57574613 3 5 Address:{String: Valid:false} CreatedAt:2020-11-03 16:41:23.651247 +0800 CST m=+0.008839280 UpdatedAt:{Time:2020-11-03 16:41:23 +0800 CST Valid:true}}

# address has value
{ID:6 Name:Zoey, Thomas Gender:2 Phone:+500 88 74 851 2657 Address:{String:47 Washington Ave,
Derby Center, NH, 99054 Valid:true} CreatedAt:2020-11-03 16:45:22.312799 +0800 CST m=+0.009339876 UpdatedAt:{Time:2020-11-03 16:45:22 +0800 CST Valid:true}}
```

fields with fix value

```go
userFactory := factory.New(
  func() interface{} { return &User{CreatedAt: time.Now()} }, // constructor with fixed value
  attr.Seq("ID", 1, "id"),
  attr.Str("Name", randutil.NameRander(3), "name"),
  attr.Int("Gender", randutil.IntRander(1, 2), "gender"),
  attr.Str("Phone", randomdata.PhoneNumber, "phone"),
  attr.Str("Address", randomdata.Address, "address"),
  attr.Time("UpdatedAt", randutil.TimeRander(time.Now(), time.Now().Add(30*time.Hour)), "updated_at"),
).Table("users").Fix("CreatedAt", "created_at")
```