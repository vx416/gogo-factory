# gogo-factory

gogo-factory is a fixtures replacement with a clear definition syntax, which inspired by [factory-go](https://github.com/bluele/factory-go) and [factory_bot](https://github.com/thoughtbot/factory_bot).

## Install

```shell
go get -u github.com/vx416/gogo-factory
```

## Getting Started

- [Defining Factories](#defining-factories)
  - [Provide initial object](#provide-initial-object)
  - [Assign value to field](#assign-value-to-field)
  - [Customize value with other fields](#customize-value-with-other-fields)
- [Building Objects](#building-objects)
  - [Build one or many objects](#build-one-or-many-objects)
  - [Setup building context](#setup-building-context)
- [Factory Associations](#factory-associations)
  - [BelongsTo association](#belongsto-association)
  - [HasOne or HasMany association](#hasone-or-hasmany-association)
  - [ManyToMany association](#manytomany-association)


### Defining Factories

This section introduce how to define a basic fixture factory. This `Employee` struct is sample struct for this section.

```go
type Gender int8

const (
  Male Gender = iota+1
  Female
)

type Employee struct {
  ID          int64
  Name        string  
  Gender      Gender  
  Age         *int32  
  Phone       string  
  Salary      float64
  CreatedAt   time.Time
  UpdatedAt   sql.NullTime
}
```

#### Provide initial object

The first argument of factory construction is initial object. The built object from factory is always a pointer.

```go
var EmployeeFactory = factory.New(
  &Employee{},
)

employee := EmployeeFactory.MustBuild().(*Employee)
```

The initial object can be used to setup initial fields.

```go
var EmployeeFactory = factory.New(
  &Employee{Name: "david"},
)

employee := EmployeeFactory.MustBuild().(*Employee)
```

#### Assign value to field

The rest of arguments of factory construction is slice of `Attributer` which is used to assign value to object's field, there are some helper functions to construct different types of `Attributer`. The `Attributer` construction require three arguments, `fieldName`, `genFunc`, and `colName`.

- fieldName: Assigned field's name

- genFunc: Value generated function

- colName: Column name in database

`types of Attributes: Int, String, Uint, Time, Bytes, Bool Inteface`

```go

var (
  id = 0
  idGenFunc = func() int {
    id++
    return id
  }
)


var EmployeeFactory = factory.New(
  &Employee{},
  attr.Int("ID", idGenFunc, "id"),
  attr.Str("Name", func() string { return "test_user" }), // you can omit the colName argument, so that this field will not be considered as column
)

employee := EmployeeFactory.MustBuild().(*Employee)
```

Customizing each `genFunc` is annoying, therefore, gogo-factory provide built-in `genFunc` to generate value with different context:

##### generate random value

```go
var EmployeeFactory = factory.New(
  &Employee{},
  attr.Int("ID", genutil.RandInt(1, 100)), // generate random int value between 1 to 100
  attr.Str("Name", genutil.RandName(3)), // generate name with gender flag (1:male, 2:female, 3:random)
  attr.Float("Salary", genutil.RandFloatSet(6.8, 5.2, 100.5, 10.5)), // generate random float value from given float set
  attr.Time("CreatedAt", genutil.RandTime(time.Now().Add(-10*24*time.Hour), time.Now())), // generate random time value between now()-10days to now()
)

employee := EmployeeFactory.MustBuild().(*Employee)
```

##### generate sequential value

```go
var EmployeeFactory = factory.New(
  &Employee{},
  attr.Int("ID", genutil.SeqInt(1, 1)), // generate sequential int value start by 1 and delta 1
  attr.Str("Name", genutil.SeqStrSet({"david", "vic", "shelly")), // generate string value from given slice sequentially
  attr.Float("Salary", genutil.SeqFloat(6.8, 0.5)), // generate sequential float value start by 6.8 and delta 0.5
  attr.Time("CreatedAt", genutil.SeqTime(time.Now(), 30*time.Second)), // generate sequential time value start by now() and delta 30sec
)

employee := EmployeeFactory.MustBuild().(*Employee)
```

##### other generated functions

```go
var EmployeeFactory = factory.New(
  &Employee{},
  attr.Int("ID", genutil.SeqInt(1, 1)),
  attr.Str("Name", genutil.FixInt("vic")), // generate fixed value
  attr.Float("Salary", genutil.SeqFloat(6.8, 0.5)),
  attr.Time("CreatedAt", genutil.Now()), // generate current time
  attr.Time("UpdatedAt", genutil.Now()), // support Scanner interface
  attr.Int("Age", genutil.SeqInt(18, 3)), // support pointer value
)

employee := EmployeeFactory.MustBuild().(*Employee)
```

#### Customize value with other fields

In some context, a value of fields which combined by other field(e.g ID) can debugger easier. The `Attributer` interface provide the `Process` method to process field.

```go
var EmployeeFactory = factory.New(
  &Employee{},
  attr.Int("ID", genutil.SeqInt(1, 1)),
  attr.Str("Name", genutil.FixInt("vic")).Process(func(a attr.Attributer) error{
    e := a.GetObject().(*Employee)
    name := a.GetVal().(string)
    return a.SetVal(fmt.Sprintf("%s-%d", name, e.ID))
  }),
)

employee := EmployeeFactory.MustBuild().(*Employee)
```

### Building Objects

#### Build one or many objects

gogo-factory can allow user to build one or many objects in one time, moreover, you can use `Insert` method to build a object meanwhile insert this object into database.

```go
factory.Opt().SetDB(db, "sqlite3") // setup global sql.DB instance and database type

var EmployeeFactory = factory.New(
  &Employee{},
  attr.Int("ID", genutil.SeqInt(1, 1)),
  attr.Str("Name", genutil.SeqStrSet({"david", "vic", "shelly")),
  attr.Float("Salary", genutil.SeqFloat(6.8, 0.5)),
  attr.Time("CreatedAt", genutil.SeqTime(time.Now(), 30*time.Second)),
).Table("employees") // setup table name

// build a object
employeeData, err := EmployeeFactory.Build()
if err == nil{
  employee := employeeData.(*Employee)
}
employee := EmployeeFactory.MustBuild().(*Employee) // panic if error not nil

// build many objects
employeesData, err := EmployeeFactory.BuildN(10)
if err == nil{
  employees := employeeData.([]*Employee)
}
employees := EmployeeFactory.MustBuildN(10).([]*Employee) // panic if error not nil


// build and insert object
employee := EmployeeFactory.MustInsert().(*Employee) // panic if error not nil
employees := EmployeeFactory.MustInsertN(10).([]*Employee) // panic if error not nil
```

#### Setup building context

When you invoke factory's method, factory will return cloned factory object which wont affect old factory building context.

#### omit some attributes

```go
// old factory object
employeeFactory := factory.New(
  &Employee{},
  attr.Int("ID", genutil.SeqInt(1, 1)),
  attr.Int("Gender", genutil.RandIntSet(1, 2)),
  attr.Time("UpdatedAt", genutil.RandTime(minTime, maxTime)),
)

// got cloned factory object
newFactory := employeeFactory.Omit("Gender")
e1 := employeeFactory.MustBuild().(*Employee)
e2 := newFactory.MustBuild().(*Employee)
assert.Equal(t, e1.ID, int64(1))
assert.NotZero(t, e1.Gender)

assert.Equal(t, e2.ID, int64(2))
assert.Zero(t, e2.Gender)
```

#### overwrite or append attributes

```go
// old factory object
employeeFactory := factory.New(
  &Employee{},
  attr.Int("ID", genutil.SeqInt(1, 1)),
  attr.Int("Gender", genutil.RandIntSet(1, 2)),
)

// got cloned factory object
newFactory := employeeFactory.Attrs(
  attr.Int("Gender", genutil.FixInt(3)),
  attr.Int("Age", genutil.SeqInt(20, 5)),
)

e1 := employeeFactory.MustBuild().(*Employee)
assert.Equal(t, e1.ID, int64(1))
assert.Nil(t, e1.Age)
assert.Less(t, int(e1.Gender), 3)

e2 := newFactory.MustBuild().(*Employee)
assert.Equal(t, e2.ID, int64(2))
assert.NotNil(t, e2.Age)
assert.Equal(t, e2.Gender, Gender(3))
```

### Factory Associations

gogo-factory support association between factories. You can combine objects and insert data across tables one time by building the factory's association.  

Here are sample structs:

```go
// domain <- 1:n -> specialty <- n:1- > employee <- m:m -> project
type Employee struct {
  ID                int64        `db:"id" gorm:"column:id"`
  Name              string       `db:"name" gorm:"column:name"`
  Gender            Gender       `db:"gender" gorm:"column:gender"`
  Age               *int32       `db:"age" gorm:"column:age"`
  Phone             string       `db:"phone" gorm:"column:phone"`
  Salary            float64      `db:"salary" gorm:"column:salary"`
  Specialty         *Specialty   `gorm:"-"`
  SecondSpecialties []*Specialty `gorm:"-"`
  Projects          []*Project   `gorm:"-"`
  CreatedAt         time.Time    `db:"created_at" gorm:"column:created_at"`
  UpdatedAt         sql.NullTime `db:"updated_at" gorm:"column:updated_at"`
}

type Project struct {
  ID        int64       `db:"id" gorm:"column:id"`
  Name      string      `db:"name" gorm:"column:name"`
  Employees []*Employee `gorm:"-"`
  Tasks     []*Task     `gorm:"-"`
  Deadline  time.Time   `db:"deadline" gorm:"column:deadline"`
}

type Task struct {
  ID        int64     `db:"id" gorm:"column:id"`
  Name      string    `db:"name" gorm:"column:name"`
  ProjectID int64     `db:"project_id" gorm:"column:project_id"`
  Deadline  time.Time `db:"deadline" gorm:"column:deadline"`
}

type Specialty struct {
  ID       int64         `db:"id" gorm:"column:id"`
  Name     string        `db:"name" gorm:"column:name"`
  OwnerID  sql.NullInt64 `db:"owner_id" gorm:"column:owner_id"`
  Domain   *Domain
  DomainID sql.NullInt64 `db:"domain_id" gorm:"column:domain_id"`
}

type Domain struct {
  ID   int64  `db:"id" gorm:"column:id"`
  Name string `db:"name" gorm:"column:name"`
}
```

#### BelongsTo association

To build association, we need to convert `Factory` object to `Association` object. This `Association` object is used to store association context information:

- ReferField: the field referred by the foreign key
  - In the example, the specialties table has a foreign key (domain_id), which refers to the ID field of Domain struct.
- ForeignKey: the name of foreign key in the table
  - In the example, the foreign key is domain_id in specialties table.
- ForeignField: the field which represents the foreign key will be assigned by the ReferField
  - note: the ForeignField can be empty


```go
var DomainFactory = factory.New(
  &Domain{},
  attr.Int("ID", genutil.SeqInt(1, 1), "id"),
  attr.Str("Name", genutil.RandAlph(10), "name"),
).Table("domains")

type SpecialtyExt struct {
  *factory.Factory
}

func (f *SpecialtyExt) BelongsToDomain(domainFactory *factory.Factory) *SpecialtyExt {
  domainAss := domainFactory.ToAssociation().ReferField("ID").ForeignKey("domain_id").ForeignField("DomainID")
  return &SpecialtyExt{f.BelongsTo("Domain", domainAss)} // use belongs to method will cloned a new factory
}

var SpecialtyFactory = &SpecialtyExt{
  factory.New(
    &Specialty{},
    attr.Int("ID", genutil.SeqInt(1, 1), "id"),
    attr.Str("Name", genutil.RandStrSet("design", "programming", "analysis", "management"), "name"),
  ).Table("specialties"),
}

spec := SpecialtyFactory.BelongsToDomain(DomainFactory).MustInsert().(*Specialty)
suite.Assert().NotZero(spec.ID)
suite.Assert().NotNil(spec.Domain)
suite.Assert().NotZero(spec.Domain.ID)
suite.Equal(spec.DomainID.Int64, spec.Domain.ID)
suite.Assert().NotEmpty(spec.Domain.Name)

specs := SpecialtyFactory.BelongsToDomain(DomainFactory).MustInsertN(5).([]*Specialty)
suite.Len(specs, 5)
```

```sql
-- MustInsert
INSERT INTO domains (id, name) VALUES (?, ?) [1 IPDZicTGXD]
INSERT INTO specialties (id, name, domain_id) VALUES (?, ?, ?) [1 analysis 1]
-- MustInsertN(5)
INSERT INTO domains (id, name) VALUES (?, ?) [2 oTMyziTclv]
INSERT INTO specialties (domain_id, id, name) VALUES (?, ?, ?) [2 2 programming]
INSERT INTO domains (id, name) VALUES (?, ?) [3 ATUrWrnfdY]
INSERT INTO specialties (name, domain_id, id) VALUES (?, ?, ?) [design 3 3]
INSERT INTO domains (name, id) VALUES (?, ?) [gOMxfdPSTz 4]
INSERT INTO specialties (domain_id, name, id) VALUES (?, ?, ?) [4 management 4]
INSERT INTO domains (id, name) VALUES (?, ?) [5 NsmgDRwbvP]
INSERT INTO specialties (name, id, domain_id) VALUES (?, ?, ?) [management 5 5]
INSERT INTO domains (id, name) VALUES (?, ?) [6 cdevwdAyKd]
INSERT INTO specialties (name, id, domain_id) VALUES (?, ?, ?) [analysis 6 6]
```

#### HasOne or HasMany association

HasOne or HasMany association information context:

- ReferField: In the example, the `employees` table is the referred table whose referred column is id and field is ID.
- ForeignKey: In the example, the `specialties` table has foreign key owner_id which refer to `employees` id.
- ForeignField: In the example, the `Specialty` struct has OwnerID field which corresponds to foreign key (owner_id).

```go
type EmployeeExt struct {
  *factory.Factory
}

func (f *EmployeeExt) HasOneSpecialty(spec *SpecialtyExt) *EmployeeExt {
  specAss := spec.ToAssociation().ReferField("ID").ForeignField("OwnerID").ForeignKey("owner_id")
  return &EmployeeExt{f.HasOne("Specialty", specAss)}
}

func (f *EmployeeExt) HasManySecondSpecialties(spec *SpecialtyExt, num int32) *EmployeeExt {
  specAss := spec.ToAssociation().ReferField("ID").ForeignField("OwnerID").ForeignKey("owner_id")
  return &EmployeeExt{f.HasMany("SecondSpecialties", specAss, num)}
}

var EmployeeFactory = &EmployeeExt{
  factory.New(
    &Employee{},
    attr.Int("ID", genutil.SeqInt(1, 1), "id"),
    attr.Str("Name", genutil.RandName(3), "name"),
    attr.Int("Gender", genutil.RandInt(1, 2), "gender"),
    attr.Int("Age", genutil.RandInt(18, 60), "age"),
    attr.Float("Salary", genutil.RandFloat(6.5, 12.8), "salary"),
    attr.Str("Phone", genutil.RandAlph(10), "phone"),
    attr.Time("CreatedAt", genutil.Now(), "created_at"),
  ).Table("employees"),
}

employee := EmployeeFactory.HasOneSpecialty(SpecialtyFactory).
  HasManySecondSpecialties(SpecialtyFactory, 5).MustInsert().(*Employee)
suite.Assert().NotZero(employee.ID)
suite.Assert().NotNil(employee.Specialty)
suite.Assert().Len(employee.SecondSpecialties, 5)
```

```sql
INSERT INTO employees (age, salary, phone, created_at, id, name, gender) VALUES (?, ?, ?, ?, ?, ?, ?) [32 11.275899307789922 psuhTzivdX 2020-11-15 14:20:40.076026 +0800 CST m=+0.079557615 1 Joshua, Robinson 1]
INSERT INTO specialties (owner_id, id, name) VALUES (?, ?, ?) [{8 true} 1 programming]
INSERT INTO specialties (id, name, owner_id) VALUES (?, ?, ?) [2 design {8 true}]
INSERT INTO specialties (id, name, owner_id) VALUES (?, ?, ?) [3 analysis {8 true}]
INSERT INTO specialties (owner_id, id, name) VALUES (?, ?, ?) [{8 true} 4 management]
INSERT INTO specialties (id, name, owner_id) VALUES (?, ?, ?) [5 programming {8 true}]
INSERT INTO specialties (owner_id, id, name) VALUES (?, ?, ?) [{8 true} 6 design]
```

The factory association can be connected, for example, `Employee` has many `Specialty`, and `Specialty` belongs to `Domain`.

```go
specAndDomainFactory := SpecialtyFactory.BelongsToDomain(DomainFactory)
employee := EmployeeFactory.HasOneSpecialty(specAndDomainFactory).
  HasManySecondSpecialties(specAndDomainFactory, 5).MustInsert().(*Employee)
suite.Assert().NotZero(employee.ID)
suite.Assert().NotNil(employee.Specialty)
suite.Assert().NotNil(employee.Specialty.Domain)
```

```sql
INSERT INTO employees (name, gender, age, salary, phone, created_at, id) VALUES (?, ?, ?, ?, ?, ?, ?) [Joshua, Jackson 1 36 12.649925079901706 kyAOLjsSxk 2020-11-15 14:25:33.125297 +0800 CST m=+0.017510176 1]
INSERT INTO domains (id, name) VALUES (?, ?) [1 dNPiJWShix]
INSERT INTO specialties (owner_id, domain_id, id, name) VALUES (?, ?, ?, ?) [{1 true} 1 1 analysis]
INSERT INTO domains (name, id) VALUES (?, ?) [RWkktbtcwx 2]
INSERT INTO specialties (name, owner_id, domain_id, id) VALUES (?, ?, ?, ?) [programming {1 true} 2 2]
INSERT INTO domains (name, id) VALUES (?, ?) [rKfZiSNxLW 3]
INSERT INTO specialties (id, name, owner_id, domain_id) VALUES (?, ?, ?, ?) [3 analysis {1 true} 3]
INSERT INTO domains (name, id) VALUES (?, ?) [VQKSqGNrKA 4]
INSERT INTO specialties (id, name, owner_id, domain_id) VALUES (?, ?, ?, ?) [4 programming {1 true} 4]
INSERT INTO domains (name, id) VALUES (?, ?) [LnwDycoyXo 5]
INSERT INTO specialties (domain_id, id, name, owner_id) VALUES (?, ?, ?, ?) [5 5 programming {1 true}]
INSERT INTO domains (name, id) VALUES (?, ?) [DTHQGwCLXz 6]
INSERT INTO specialties (owner_id, id, name, domain_id) VALUES (?, ?, ?, ?) [{1 true} 6 analysis 6]
```

#### ManyToMany association

ManyToMany association is more complex, so it need more information:

- JoinTable: The join table represent a table which link two tables. You can use `Attributer` to set column value (e.g id) for join table.
  - In the example, the join table is `employees_projects` which has `employee_id` and `project_id` as foreign keys.
- ReferField: The field of current struct has value of foreign key in the join table.
  - In the example, the current struct `Employee` has `ID` field which corresponds to `employee_id` in join table.
- ReferColumn: The column name of join table which corresponds to `ReferField`.
  - In the example, the `ReferColumn` will be `employee_id`.
- ForeignField: The field of associated struct has value of foreign key in the join table.
  - In the example, the associated struct `Project` has `ID` field which corresponds to `project_id` in join table.
- ForeignKey: The column name of join table which corresponds to `ForeignField`.
  - In the example, the `ReferColumn` will be `project_id`.
- AssociatedField: The field of associated struct `Project` whose element is current struct `Employee`.
  - In the example, the `AssociatedField` will be `Employees`.
  - note: `AssociatedField` should be slice.

```go
type EmployeeExt struct {
  *factory.Factory
}

func (f *EmployeeExt) HasManyProjects(prj *ProjectExt, num int32) *EmployeeExt {
  prjAss := prj.ToAssociation().ReferField("ID").ReferColumn("employee_id").
    ForeignField("ID").ForeignKey("project_id").AssociatedField("Employees").
    JoinTable("employees_projects", attr.Int("ID", genutil.SeqInt(1, 1), "id"))
  return &EmployeeExt{f.ManyToMany("Projects", prjAss, num)}
}


type ProjectExt struct {
  *factory.Factory
}

var ProjectFactory = &ProjectExt{
  factory.New(
    &Project{},
    attr.Int("ID", genutil.SeqInt(1, 1), "id")
    attr.Str("Name", genutil.RandUUID(), "name"),
    attr.Time("Deadline", genutil.SeqTime(time.Now(), 100*time.Hour), "deadline"),
  ).Table("projects"),
}


employee := EmployeeFactory.HasManyProjects(ProjectFactory, 2).MustInsert().(*Employee)
suite.Assert().NotZero(employee.ID)
suite.Assert().Len(employee.Projects, 2)
suite.Assert().Len(employee.Projects[0].Employees, 1)
suite.Assert().Len(employee.Projects[1].Employees, 1)
suite.Assert().Equal(employee.Projects[1].Employees[0].ID, employee.ID)
```

```sql
INSERT INTO employees (created_at, id, name, gender, age, salary, phone) VALUES (?, ?, ?, ?, ?, ?, ?) [2020-11-15 14:59:22.99355 +0800 CST m=+0.015319184 1 Zoey, Thomas 1 36 10.318368751716168 dUHjJKjLoV]
INSERT INTO projects (name, deadline, id) VALUES (?, ?, ?) [bb19d463-dabf-4afe-88e2-296d44a5f09f 2020-11-15 14:59:22.981192 +0800 CST m=+0.002961345 1]
INSERT INTO projects (id, name, deadline) VALUES (?, ?, ?) [2 a4dafce3-8e4b-4f34-a56c-2d30c79fcfb4 2020-11-19 18:59:22.981192 +0800 CST m=+360000.002961345]
INSERT INTO employees_projects (employee_id, project_id, id) VALUES (?, ?, ?) [1 1 1]
INSERT INTO employees_projects (employee_id, project_id, id) VALUES (?, ?, ?) [1 2 2]
```

Connect all association together.

```go

type ProjectExt struct {
  *factory.Factory
}

func (f *ProjectExt) HasManyTasks(taskFactory *factory.Factory, num int32) *ProjectExt {
  taskAss := taskFactory.ToAssociation().ReferField("ID").ForeignKey("project_id").ForeignField("ProjectID")
  return &ProjectExt{f.HasMany("Tasks", taskAss, num)}
}

var TaskFactory = factory.New(
  &Task{},
  attr.Int("ID", genutil.SeqInt(1, 1), "id")
  attr.Str("Name", genutil.RandUUID(), "name"),
  attr.Time("Deadline", genutil.SeqTime(time.Now(), 100*time.Hour), "deadline"),
).Table("tasks")

spec := SpecialtyFactory.BelongsToDomain(DomainFactory)
proj := ProjectFactory.HasManyTasks(TaskFactory, 10)
employee := EmployeeFactory.HasOneSpecialty(spec).HasManySecondSpecialties(spec, 3).
  HasManyProjects(proj, 2).MustInsert().(*Employee)
suite.Assert().NotZero(employee.ID)
suite.Assert().Len(employee.Projects, 2)
suite.Assert().Equal(employee.Projects[1].Employees[0].ID, employee.ID)
suite.Assert().NotNil(employee.Specialty)
suite.Assert().NotNil(employee.Specialty.Domain)
suite.Assert().Len(employee.Projects[0].Tasks, 10)
```

```sql
INSERT INTO employees (name, gender, age, salary, phone, created_at, id) VALUES (?, ?, ?, ?, ?, ?, ?) [Isabella, Williams 2 54 12.323119013120419 VrfhHyrTAL 2020-11-15 15:35:45.777974 +0800 CST m=+0.014762944 1]
INSERT INTO domains (id, name) VALUES (?, ?) [1 aLvKjKUdXY]
INSERT INTO specialties (id, name, owner_id, domain_id) VALUES (?, ?, ?, ?) [1 design {1 true} 1]
INSERT INTO domains (id, name) VALUES (?, ?) [2 fCLOddJBBF]
INSERT INTO specialties (id, name, owner_id, domain_id) VALUES (?, ?, ?, ?) [2 design {1 true} 2]
INSERT INTO domains (id, name) VALUES (?, ?) [3 riivdoxaxJ]
INSERT INTO specialties (id, name, owner_id, domain_id) VALUES (?, ?, ?, ?) [3 programming {1 true} 3]
INSERT INTO domains (id, name) VALUES (?, ?) [4 QChCzKWQvs]
INSERT INTO specialties (id, name, owner_id, domain_id) VALUES (?, ?, ?, ?) [4 programming {1 true} 4]
INSERT INTO projects (deadline, id, name) VALUES (?, ?, ?) [2020-11-15 15:35:45.766096 +0800 CST m=+0.002885016 1 2f4a8e39-66d1-4cd4-a0cb-96d892a822a5]
INSERT INTO tasks (name, deadline, id, project_id) VALUES (?, ?, ?, ?) [cfa26a28-0e3d-48c0-9e6f-582d2abbb9be 2020-11-15 15:35:45.766097 +0800 CST m=+0.002886438 1 1]
INSERT INTO tasks (id, project_id, name, deadline) VALUES (?, ?, ?, ?) [2 1 c8b323c4-f66d-4071-be73-0409d04eeccd 2020-11-19 19:35:45.766097 +0800 CST m=+360000.002886438]
INSERT INTO tasks (project_id, name, deadline, id) VALUES (?, ?, ?, ?) [1 9446c004-f5c0-43f9-a8d0-5864cc86b357 2020-11-23 23:35:45.766097 +0800 CST m=+720000.002886438 3]
INSERT INTO tasks (name, deadline, id, project_id) VALUES (?, ?, ?, ?) [2b57a035-548f-4269-9cd9-001f2f523bbf 2020-11-28 03:35:45.766097 +0800 CST m=+1080000.002886438 4 1]
INSERT INTO tasks (deadline, id, project_id, name) VALUES (?, ?, ?, ?) [2020-12-02 07:35:45.766097 +0800 CST m=+1440000.002886438 5 1 4ee7b87e-430c-488b-b991-0918b1c59b12]
INSERT INTO tasks (project_id, name, deadline, id) VALUES (?, ?, ?, ?) [1 6f9b8921-1d0b-407b-85c9-08a6ce3a192a 2020-12-06 11:35:45.766097 +0800 CST m=+1800000.002886438 6]
INSERT INTO tasks (deadline, id, project_id, name) VALUES (?, ?, ?, ?) [2020-12-10 15:35:45.766097 +0800 CST m=+2160000.002886438 7 1 803f1d9e-66f5-4dfe-96d5-6bb385ef457c]
INSERT INTO tasks (name, deadline, id, project_id) VALUES (?, ?, ?, ?) [0b380cc1-513f-4c5e-b627-62d88b463ab5 2020-12-14 19:35:45.766097 +0800 CST m=+2520000.002886438 8 1]
INSERT INTO tasks (deadline, id, project_id, name) VALUES (?, ?, ?, ?) [2020-12-18 23:35:45.766097 +0800 CST m=+2880000.002886438 9 1 d9700d31-7c98-49d8-a0e9-1643d0b10fcc]
INSERT INTO tasks (project_id, name, deadline, id) VALUES (?, ?, ?, ?) [1 7f09ca4c-9bbd-4188-85d4-17db84097bbd 2020-12-23 03:35:45.766097 +0800 CST m=+3240000.002886438 10]
INSERT INTO projects (deadline, id, name) VALUES (?, ?, ?) [2020-11-19 19:35:45.766096 +0800 CST m=+360000.002885016 2 c014be25-9026-4d57-99e2-4d0d97245efa]
INSERT INTO tasks (id, project_id, name, deadline) VALUES (?, ?, ?, ?) [11 2 c578d2fc-03c2-48f7-b08b-e96fe09531a3 2020-12-27 07:35:45.766097 +0800 CST m=+3600000.002886438]
INSERT INTO tasks (id, project_id, name, deadline) VALUES (?, ?, ?, ?) [12 2 3eeafd4d-931e-4b59-bc9b-6e7cfd610bd9 2020-12-31 11:35:45.766097 +0800 CST m=+3960000.002886438]
INSERT INTO tasks (id, project_id, name, deadline) VALUES (?, ?, ?, ?) [13 2 74470827-2cda-4c30-a963-9a1d0bb4e253 2021-01-04 15:35:45.766097 +0800 CST m=+4320000.002886438]
INSERT INTO tasks (deadline, id, project_id, name) VALUES (?, ?, ?, ?) [2021-01-08 19:35:45.766097 +0800 CST m=+4680000.002886438 14 2 3e42ed58-1eee-4da4-a092-97c038c3a871]
INSERT INTO tasks (name, deadline, id, project_id) VALUES (?, ?, ?, ?) [4157556e-2fa6-448c-86ef-6ffe29789ddf 2021-01-12 23:35:45.766097 +0800 CST m=+5040000.002886438 15 2]
INSERT INTO tasks (id, project_id, name, deadline) VALUES (?, ?, ?, ?) [16 2 7d4e68ba-2037-41c1-a9b1-2f113aa701c3 2021-01-17 03:35:45.766097 +0800 CST m=+5400000.002886438]
INSERT INTO tasks (id, project_id, name, deadline) VALUES (?, ?, ?, ?) [17 2 034e05f9-b394-460f-a467-7ce9374b41e5 2021-01-21 07:35:45.766097 +0800 CST m=+5760000.002886438]
INSERT INTO tasks (project_id, name, deadline, id) VALUES (?, ?, ?, ?) [2 3675914a-09b5-423c-b575-77b109ee14ba 2021-01-25 11:35:45.766097 +0800 CST m=+6120000.002886438 18]
INSERT INTO tasks (deadline, id, project_id, name) VALUES (?, ?, ?, ?) [2021-01-29 15:35:45.766097 +0800 CST m=+6480000.002886438 19 2 2dc6f22e-8a79-4172-881a-ebf30eaf51a0]
INSERT INTO tasks (deadline, id, project_id, name) VALUES (?, ?, ?, ?) [2021-02-02 19:35:45.766097 +0800 CST m=+6840000.002886438 20 2 35eb35d7-0e23-4871-962a-dce0280388ae]
INSERT INTO employees_projects (employee_id, project_id, id) VALUES (?, ?, ?) [1 1 1]
INSERT INTO employees_projects (employee_id, project_id, id) VALUES (?, ?, ?) [1 2 2]
```
