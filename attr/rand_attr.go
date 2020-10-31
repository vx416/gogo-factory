package attr

import (
	"fmt"
	"math/rand"
	"time"
)

// RandInt create int attributer with rand number in range of min and max
func RandInt(name string, min, max int, options ...string) Attributer {
	return &randIntAttr{
		name:    name,
		colName: getColName(options),
		min:     min,
		max:     max,
	}
}

type randIntAttr struct {
	name    string
	colName string
	min     int
	max     int
	val     int
	process Processor
}

func (attr *randIntAttr) Process(procFunc Processor) Attributer {
	attr.process = procFunc
	return attr
}

func (attr randIntAttr) GetVal() interface{} {
	return attr.val
}

func (attr *randIntAttr) SetVal(val interface{}) error {
	realVal, ok := val.(int)
	if !ok {
		return fmt.Errorf("set attribute val: val %+v is not int", val)
	}

	attr.val = realVal
	return nil
}

func (attr randIntAttr) ColName() string {
	return attr.colName
}

func (attr *randIntAttr) Gen(data interface{}) (interface{}, error) {
	attr.val = randIntIn(attr.min, attr.max)
	if attr.process != nil {
		if err := attr.process(attr, data); err != nil {
			return nil, err
		}
	}

	return attr.val, nil
}

func (randIntAttr) Kind() Type {
	return IntAttr
}

func (attr randIntAttr) Name() string {
	return attr.name
}

// RandUint create uint attributer with rand number in range of min and max
func RandUint(name string, min, max uint, options ...string) Attributer {
	return &randUintAttr{
		name:    name,
		colName: getColName(options),
		min:     min,
		max:     max,
	}
}

type randUintAttr struct {
	min     uint
	max     uint
	val     uint
	name    string
	colName string
	process Processor
}

func (attr *randUintAttr) Process(procFunc Processor) Attributer {
	attr.process = procFunc
	return attr
}

func (attr randUintAttr) GetVal() interface{} {
	return attr.val
}

func (attr *randUintAttr) SetVal(val interface{}) error {
	realVal, ok := val.(uint)
	if !ok {
		return fmt.Errorf("set attribute val: val %+v is not uint", val)
	}

	attr.val = realVal
	return nil
}

func (attr randUintAttr) ColName() string {
	return attr.colName
}

func (attr *randUintAttr) Gen(data interface{}) (interface{}, error) {
	attr.val = randUintIn(attr.min, attr.max)
	if attr.process != nil {
		if err := attr.process(attr, data); err != nil {
			return nil, err
		}
	}

	return attr.val, nil
}

func (randUintAttr) Kind() Type {
	return UintAttr
}

func (attr randUintAttr) Name() string {
	return attr.name
}

func randIntIn(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func randUintIn(min, max uint) uint {
	rand.Seed(time.Now().UnixNano())
	return uint(rand.Intn(int(max)-int(min)) + int(min))
}

// RandFloat create float attributer with rand number in range of min and max
func RandFloat(name string, min, max float64, options ...string) Attributer {
	return &randFloatAttr{
		name:    name,
		colName: getColName(options),
		min:     min,
		max:     max,
	}
}

type randFloatAttr struct {
	min     float64
	max     float64
	val     float64
	name    string
	colName string
	process Processor
}

func (attr *randFloatAttr) Process(procFunc Processor) Attributer {
	attr.process = procFunc
	return attr
}

func (attr randFloatAttr) GetVal() interface{} {
	return attr.val
}

func (attr *randFloatAttr) SetVal(val interface{}) error {
	realVal, ok := val.(float64)
	if !ok {
		return fmt.Errorf("set attribute val: val %+v is not float64", val)
	}

	attr.val = realVal
	return nil
}

func (attr randFloatAttr) ColName() string {
	return attr.colName
}

func (attr *randFloatAttr) Gen(data interface{}) (interface{}, error) {
	attr.val = randFloat64(attr.min, attr.max)
	if attr.process != nil {
		if err := attr.process(attr, data); err != nil {
			return nil, err
		}
	}
	return attr.val, nil
}

func (randFloatAttr) Kind() Type {
	return FloatAttr
}

func (attr randFloatAttr) Name() string {
	return attr.name
}

func randFloat64(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	minInt, maxInt := int(min), int(max)
	randInt := randIntIn(minInt, maxInt)
	return float64(randInt) + rand.Float64()
}

// RandStr create string attributer with rand string within string slice
func RandStr(name string, randSet []string, options ...string) Attributer {
	return &randStrAttr{
		name:     name,
		colName:  getColName(options),
		randSet:  randSet,
		maxIndex: len(randSet) - 1,
	}
}

type randStrAttr struct {
	randSet  []string
	maxIndex int
	val      string
	name     string
	colName  string
	process  Processor
}

func (attr *randStrAttr) Process(procFunc Processor) Attributer {
	attr.process = procFunc
	return attr
}

func (attr randStrAttr) GetVal() interface{} {
	return attr.val
}

func (attr *randStrAttr) SetVal(val interface{}) error {
	realVal, ok := val.(string)
	if !ok {
		return fmt.Errorf("set attribute val: val %+v is not string", val)
	}

	attr.val = realVal
	return nil
}

func (attr randStrAttr) ColName() string {
	return attr.colName
}

func (attr *randStrAttr) Gen(data interface{}) (interface{}, error) {
	index := randIntIn(0, attr.maxIndex)
	attr.val = attr.randSet[index]
	if attr.process != nil {
		if err := attr.process(attr, data); err != nil {
			return nil, err
		}
	}
	return attr.val, nil
}

func (randStrAttr) Kind() Type {
	return StringAttr
}

func (attr randStrAttr) Name() string {
	return attr.name
}

// RandTime create time attributer with rand time in range of min and max
func RandTime(name string, min, max time.Time, options ...string) Attributer {
	return &randTimeAttr{
		name:    name,
		colName: getColName(options),
		minTime: int(min.Unix()),
		maxTime: int(max.Unix()),
	}
}

type randTimeAttr struct {
	minTime int
	maxTime int
	val     time.Time
	name    string
	colName string
	process Processor
}

func (attr *randTimeAttr) Process(procFunc Processor) Attributer {
	attr.process = procFunc
	return attr
}

func (attr randTimeAttr) GetVal() interface{} {
	return attr.val
}

func (attr *randTimeAttr) SetVal(val interface{}) error {
	realVal, ok := val.(time.Time)
	if !ok {
		return fmt.Errorf("set attribute val: val %+v is not time.Time", val)
	}

	attr.val = realVal
	return nil
}

func (attr randTimeAttr) ColName() string {
	return attr.colName
}

func (attr *randTimeAttr) Gen(data interface{}) (interface{}, error) {
	timeUnix := randIntIn(attr.minTime, attr.maxTime)
	attr.val = time.Unix(int64(timeUnix), 0)
	if attr.process != nil {
		if err := attr.process(attr, data); err != nil {
			return nil, err
		}
	}
	return attr.val, nil
}

func (randTimeAttr) Kind() Type {
	return TimeAttr
}

func (attr randTimeAttr) Name() string {
	return attr.name
}

// RandBool create boolean attributer with random true/false value
//   the ratio parameter define the percentage of occurrence of true value
func RandBool(name string, ratio float64, options ...string) Attributer {
	return &randBoolAttr{
		name:    name,
		ratio:   ratio,
		colName: getColName(options),
	}
}

type randBoolAttr struct {
	val     bool
	ratio   float64
	name    string
	colName string
	process Processor
}

func (attr *randBoolAttr) Process(procFunc Processor) Attributer {
	attr.process = procFunc
	return attr
}

func (attr randBoolAttr) GetVal() interface{} {
	return attr.val
}

func (attr *randBoolAttr) SetVal(val interface{}) error {
	realVal, ok := val.(bool)
	if !ok {
		return fmt.Errorf("set attribute val: val %+v is not bool", val)
	}

	attr.val = realVal
	return nil
}

func (attr randBoolAttr) ColName() string {
	return attr.colName
}

func (randBoolAttr) Kind() Type {
	return BoolAttr
}

func (attr *randBoolAttr) Gen(data interface{}) (interface{}, error) {
	rand.Seed(time.Now().UnixNano())
	if rand.Float64() > attr.ratio {
		attr.val = false
	} else {
		attr.val = true
	}

	if attr.process != nil {
		if err := attr.process(attr, data); err != nil {
			return nil, err
		}
	}
	return attr.val, nil
}

func (attr randBoolAttr) Name() string {
	return attr.name
}
