package test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/suite"
	factory "github.com/vx416/gogo-factory"
)

func TestSqlite(t *testing.T) {
	db, err := initSqliteDB()
	if err != nil {
		t.Fatalf("db init failed, err:%+v", err)
	}
	s := &insertSuite{
		db:     db,
		dbType: "sqlite3",
	}
	suite.Run(t, s)
}

func employeeEq(suite *insertSuite, a *Employee, b *Employee) {
	suite.Equal(a.ID, b.ID)
	suite.Equal(a.Gender, b.Gender)
	suite.Equal(a.Salary, b.Salary)
	suite.Equal(a.Phone, b.Phone)
	if a.Age != nil {
		suite.Equal(*a.Age, *b.Age)
	}
	suite.True(a.CreatedAt.Equal(b.CreatedAt))
}

func specEq(suite *insertSuite, a *Specialty, b *Specialty) {
	suite.Equal(a.ID, b.ID)
	suite.Equal(a.Name, b.Name)
	if a.Domain != nil {
		suite.Equal(a.DomainID.Int64, b.DomainID.Int64)
		suite.Equal(a.Domain.ID, b.DomainID.Int64)
	}
	if a.OwnerID.Valid {
		suite.Equal(a.OwnerID.Int64, b.OwnerID.Int64)
	}
}

func domainEq(suite *insertSuite, a *Domain, b *Domain) {
	suite.Equal(a.ID, b.ID)
	suite.Equal(a.Name, b.Name)
}

func employeeAndSpec(suite *insertSuite, a *Employee, b *Employee, specs []*Specialty, domains []*Domain) {
	employeeEq(suite, a, b)

	specsMap := make(map[int64]*Specialty)
	for i := range specs {
		specsMap[specs[i].ID] = specs[i]
	}
	domainMap := make(map[int64]*Domain)
	if domains != nil {
		for i := range domains {
			domainMap[domains[i].ID] = domains[i]
		}
	}

	specEq(suite, a.Specialty, specsMap[a.Specialty.ID])

	for i := 0; i < len(a.SecondSpecialties); i++ {
		secondSpec := a.SecondSpecialties[i]
		specEq(suite, secondSpec, specsMap[secondSpec.ID])
	}

	if domains != nil {
		if a.Specialty.Domain != nil {
			domainEq(suite, a.Specialty.Domain, domainMap[a.Specialty.Domain.ID])
		}

		for i := 0; i < len(a.SecondSpecialties); i++ {
			secondSpec := a.SecondSpecialties[i]
			if secondSpec.Domain != nil {
				domainEq(suite, secondSpec.Domain, domainMap[secondSpec.Domain.ID])
			}
		}
	}
}

type insertSuite struct {
	db     *sql.DB
	dbType string
	suite.Suite
}

func (suite *insertSuite) SetupSuite() {
	factory.Opt().SetDB(suite.db, suite.dbType)
}

func (suite *insertSuite) TearDownTest() {
	err := Clear(suite.db)
	suite.Require().NoError(err)
}

func (suite *insertSuite) TestInsertOne() {
	employee := EmployeeFactory.MustInsert().(*Employee)
	suite.Assert().NotZero(employee.ID)
	employees, _, err := AllEmployees(suite.db, suite.dbType)
	suite.Assert().NoError(err)
	suite.Len(employees, 1)
	employeeEq(suite, employee, employees[0])

	employees2 := EmployeeFactory.MustInsertN(5).([]*Employee)
	suite.Len(employees2, 5)
	_, employeeMap, err := AllEmployees(suite.db, suite.dbType)
	suite.Assert().NoError(err)
	suite.Len(employees, 6)
	for i := 0; i < len(employees2); i++ {
		employeeEq(suite, employees2[i], employeeMap[employees2[i].ID])
	}
}

func (suite *insertSuite) TestBelongsTo() {
	spec := SpecialtyFactory.BelongsToDomain(DomainFactory).MustInsert().(*Specialty)
	suite.Assert().NotZero(spec.ID)
	suite.Assert().NotNil(spec.Domain)
	specs, err := AllSpecialties(suite.db, suite.dbType)
	suite.Require().NoError(err)
	suite.Len(specs, 1)
	specEq(suite, spec, specs[0])
	domains, err := AllDomains(suite.db, suite.dbType)
	suite.Require().NoError(err)
	suite.Len(domains, 1)
	domainEq(suite, spec.Domain, domains[0])

	specs2 := SpecialtyFactory.BelongsToDomain(DomainFactory).MustInsertN(5).([]*Specialty)
	suite.Len(specs2, 5)
	specs, err = AllSpecialties(suite.db, suite.dbType)
	suite.Require().NoError(err)
	suite.Len(specs, 6)
	for i := 1; i < len(domains); i++ {
		specEq(suite, specs2[i-1], specs[i])
	}

	domains, err = AllDomains(suite.db, suite.dbType)
	suite.Require().NoError(err)
	suite.Len(domains, 6)
	for i := 1; i < len(domains); i++ {
		domainEq(suite, specs2[i-1].Domain, domains[i])
	}
}

func (suite *insertSuite) TestHasOneAndMany() {
	employee := EmployeeFactory.HasOneSpecialty(SpecialtyFactory).
		HasManySecondSpecialties(SpecialtyFactory, 5).MustInsert().(*Employee)
	suite.Assert().NotZero(employee.ID)
	suite.Assert().NotNil(employee.Specialty)
	suite.Assert().Len(employee.SecondSpecialties, 5)

	employees, _, err := AllEmployees(suite.db, suite.dbType)
	suite.Assert().NoError(err)
	suite.Len(employees, 1)
	specs, err := AllSpecialties(suite.db, suite.dbType)
	suite.Require().NoError(err)
	suite.Len(specs, 6)
	employeeAndSpec(suite, employee, employees[0], specs, nil)
}

func (suite *insertSuite) TestHasOneAndManyN() {
	employees := EmployeeFactory.HasOneSpecialty(SpecialtyFactory).
		HasManySecondSpecialties(SpecialtyFactory, 5).MustInsertN(5).([]*Employee)
	suite.Assert().Len(employees, 5)
	suite.Assert().NotZero(employees[0].ID)
	suite.Assert().NotNil(employees[0].Specialty)
	suite.Assert().Len(employees[0].SecondSpecialties, 5)

	employees2, _, err := AllEmployees(suite.db, suite.dbType)
	suite.Assert().NoError(err)
	suite.Len(employees2, 5)
	specs, err := AllSpecialties(suite.db, suite.dbType)
	suite.Require().NoError(err)
	suite.Len(specs, 6*5)

	for i := range employees2 {
		employeeAndSpec(suite, employees[i], employees2[i], specs, nil)
	}
}

func testEmployeeSpecialtyDomain(suite *insertSuite, employees ...*Employee) {
	employees2, employeesMap, err := AllEmployees(suite.db, suite.dbType)
	suite.Assert().NoError(err)
	suite.Len(employees2, len(employees))
	specs, err := AllSpecialties(suite.db, suite.dbType)
	suite.Require().NoError(err)
	domains, err := AllDomains(suite.db, suite.dbType)
	suite.Require().NoError(err)

	for _, employee := range employees {
		suite.Assert().NotZero(employee.ID)
		suite.Assert().NotNil(employee.Specialty)
		suite.Assert().NotNil(employee.Specialty.Domain)
		suite.Assert().NotEmpty(employee.SecondSpecialties)
		for i := range employee.SecondSpecialties {
			suite.Assert().NotNil(employee.SecondSpecialties[i].Domain)
		}
		employeeAndSpec(suite, employee, employeesMap[employee.ID], specs, domains)

	}
}

func (suite *insertSuite) TestBelongsToAndHasOneAndMany() {
	specFactory := SpecialtyFactory.BelongsToDomain(DomainFactory)
	employee := EmployeeFactory.HasOneSpecialty(specFactory).
		HasManySecondSpecialties(specFactory, 5).MustInsert().(*Employee)
	testEmployeeSpecialtyDomain(suite, employee)
}

func (suite *insertSuite) TestBelongsToAndHasOneAndManyN() {
	specFactory := SpecialtyFactory.BelongsToDomain(DomainFactory)
	employees := EmployeeFactory.HasOneSpecialty(specFactory).
		HasManySecondSpecialties(specFactory, 5).MustInsertN(5).([]*Employee)
	suite.Assert().Len(employees, 5)
	suite.Assert().NotZero(employees[0].ID)
	suite.Assert().NotNil(employees[0].Specialty)
	suite.Assert().Len(employees[0].SecondSpecialties, 5)

	employees2, _, err := AllEmployees(suite.db, suite.dbType)
	suite.Assert().NoError(err)
	suite.Len(employees2, 5)
	specs, err := AllSpecialties(suite.db, suite.dbType)
	suite.Require().NoError(err)
	suite.Len(specs, 6*5)
	domains, err := AllDomains(suite.db, suite.dbType)
	suite.Require().NoError(err)
	suite.Len(domains, 6*5)
	for i := range employees2 {
		employeeAndSpec(suite, employees[i], employees2[i], specs, domains)
	}
}

func testEmployeesAndProjects(suite *insertSuite, projectNum int, employess ...*Employee) {
	employees2, _, err := AllEmployees(suite.db, suite.dbType)
	suite.NoError(err)
	employeesMap := make(map[int64]*Employee)

	for i := range employees2 {
		employeesMap[employees2[i].ID] = employees2[i]
	}

	projects2, err := AllProjects(suite.db, suite.dbType)
	suite.NoError(err)
	projectsMap := make(map[int64]*Project)
	for i := range projects2 {
		projectsMap[projects2[i].ID] = projects2[i]
	}

	employeesPrjs, err := AllEmployeesProjects(suite.db, suite.dbType)
	suite.NoError(err)
	emplyPrjsMap := make(map[int64]map[int64]*EmployeesProjects)
	for i := range employeesPrjs {
		if _, ok := emplyPrjsMap[employeesPrjs[i].EmployeeID]; !ok {
			emplyPrjsMap[employeesPrjs[i].EmployeeID] = make(map[int64]*EmployeesProjects)
		}
		emplyPrjsMap[employeesPrjs[i].EmployeeID][employeesPrjs[i].ProjectID] = employeesPrjs[i]
	}

	for _, employee := range employess {
		suite.Assert().NotZero(employee.ID)
		suite.Assert().Len(employee.Projects, projectNum)
		employee2 := employeesMap[employee.ID]
		employeeEq(suite, employee, employee2)
		for i := range employee.Projects {
			suite.Assert().Len(employee.Projects[i].Employees, 1)
			testProjectEq(suite, employee.Projects[i], projectsMap[employee.Projects[i].ID])
			_, ok := emplyPrjsMap[employee.ID][employee.Projects[i].ID]
			suite.Assert().True(ok)
		}
	}
}

func testProjectEq(suite *insertSuite, prjA *Project, prjB *Project) {
	suite.Require().NotNil(prjA)
	suite.Require().NotNil(prjB)
	suite.Assert().Equal(prjA.ID, prjB.ID)
	suite.Assert().Equal(prjA.Name, prjB.Name)
	suite.Assert().True(prjA.Deadline.Equal(prjB.Deadline))
}

func (suite *insertSuite) TestManyToMany() {
	f := EmployeeFactory.HasManyProjects(ProjectFactory, 2)
	employee := f.MustInsert().(*Employee)
	testEmployeesAndProjects(suite, 2, employee)
	employee = f.MustInsert().(*Employee)
	testEmployeesAndProjects(suite, 2, employee)
}

func (suite *insertSuite) TestManyToManyN() {
	employees := EmployeeFactory.HasManyProjects(ProjectFactory, 2).MustInsertN(5).([]*Employee)
	testEmployeesAndProjects(suite, 2, employees...)
}

func (suite *insertSuite) TestManyToManyBelongsTo() {
	spec := SpecialtyFactory.BelongsToDomain(DomainFactory)
	employee := EmployeeFactory.HasOneSpecialty(spec).HasManySecondSpecialties(spec, 3).
		HasManyProjects(ProjectFactory, 2).MustInsert().(*Employee)
	testEmployeesAndProjects(suite, 2, employee)
	testEmployeeSpecialtyDomain(suite, employee)
}

func (suite *insertSuite) TestAll() {
	spec := SpecialtyFactory.BelongsToDomain(DomainFactory)
	proj := ProjectFactory.HasManyTasks(TaskFactory, 10)
	employee := EmployeeFactory.HasOneSpecialty(spec).HasManySecondSpecialties(spec, 3).
		HasManyProjects(proj, 2).MustInsert().(*Employee)
	testEmployeesAndProjects(suite, 2, employee)
	testEmployeeSpecialtyDomain(suite, employee)
	suite.Assert().Len(employee.Projects[0].Tasks, 10)
}

func (suite *insertSuite) TestAllN() {
	spec := SpecialtyFactory.BelongsToDomain(DomainFactory)
	proj := ProjectFactory.HasManyTasks(TaskFactory, 10)
	employees := EmployeeFactory.HasOneSpecialty(spec).HasManySecondSpecialties(spec, 3).
		HasManyProjects(proj, 2).MustInsertN(3).([]*Employee)
	testEmployeesAndProjects(suite, 2, employees...)
	testEmployeeSpecialtyDomain(suite, employees...)
	for i := range employees {
		suite.Assert().Len(employees[i].Projects[0].Tasks, 10)
	}
}
