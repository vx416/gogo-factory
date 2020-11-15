package test

import (
	"time"

	factory "github.com/vx416/gogo-factory"
	"github.com/vx416/gogo-factory/attr"
	"github.com/vx416/gogo-factory/genutil"
)

var idAttr = func() attr.Attributer {
	return attr.Int("ID", genutil.SeqInt(1, 1), "id")
}

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

func (f *EmployeeExt) HasManyProjects(prj *ProjectExt, num int32) *EmployeeExt {
	prjAss := prj.ToAssociation().ReferField("ID").ReferColumn("employee_id").
		ForeignField("ID").ForeignKey("project_id").AssociatedField("Employees").
		JoinTable("employees_projects", attr.Int("ID", genutil.SeqInt(1, 1), "id"))
	return &EmployeeExt{f.ManyToMany("Projects", prjAss, num)}
}

var EmployeeFactory = &EmployeeExt{
	factory.New(
		&Employee{},
		idAttr(),
		attr.Str("Name", genutil.RandName(3), "name"),
		attr.Int("Gender", genutil.RandInt(1, 2), "gender"),
		attr.Int("Age", genutil.RandInt(18, 60), "age"),
		attr.Float("Salary", genutil.RandFloat(6.5, 12.8), "salary"),
		attr.Str("Phone", genutil.RandAlph(10), "phone"),
		attr.Time("CreatedAt", genutil.Now(), "created_at"),
	).Table("employees"),
}

var DomainFactory = factory.New(
	&Domain{},
	idAttr(),
	attr.Str("Name", genutil.RandAlph(10), "name"),
).Table("domains")

type SpecialtyExt struct {
	*factory.Factory
}

func (f *SpecialtyExt) BelongsToDomain(other *factory.Factory) *SpecialtyExt {
	domainAss := other.ToAssociation().ReferField("ID").ForeignKey("domain_id").ForeignField("DomainID")
	return &SpecialtyExt{f.BelongsTo("Domain", domainAss)}
}

var SpecialtyFactory = &SpecialtyExt{
	factory.New(
		&Specialty{},
		idAttr(),
		attr.Str("Name", genutil.RandStrSet("design", "programming", "analysis", "management"), "name"),
	).Table("specialties"),
}

type ProjectExt struct {
	*factory.Factory
}

func (f *ProjectExt) HasManyTasks(taskFactory *factory.Factory, num int32) *ProjectExt {
	taskAss := taskFactory.ToAssociation().ReferField("ID").ForeignKey("project_id").ForeignField("ProjectID")
	return &ProjectExt{f.HasMany("Tasks", taskAss, num)}
}

var ProjectFactory = &ProjectExt{
	factory.New(
		&Project{},
		idAttr(),
		attr.Str("Name", genutil.RandUUID(), "name"),
		attr.Time("Deadline", genutil.SeqTime(time.Now(), 100*time.Hour), "deadline"),
	).Table("projects"),
}

var TaskFactory = factory.New(
	&Task{},
	idAttr(),
	attr.Str("Name", genutil.RandUUID(), "name"),
	attr.Time("Deadline", genutil.SeqTime(time.Now(), 100*time.Hour), "deadline"),
).Table("tasks")
