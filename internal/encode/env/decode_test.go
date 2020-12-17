package env

import (
	"os"
	"testing"
	"time"

	"github.com/jbsmith7741/go-tools/appenderr"
	"github.com/jbsmith7741/trial"
)

/*
func TestDecoder_Unmarshal(t *testing.T) {
	type Level3 struct {
		FirstField  *string
		SecondField string `env:"second_field"`
	}

	type level2 struct {
		FirstField  *string
		SecondField string `env:"second_field"`

		privateField string // Should not be populated.

		Level3        Level3  `env:"LEVEL3"` // Should not matter if struct type is public or private. Only the field name.
		Level3Pointer *Level3 // Should initialize Level3 and store the pointer type and not panic.
	}

	type aField struct {
		Field1 string
	}

	type withprefix struct {
		WithPrefix aField
	}

	type omitprefix struct {
		NoPrefix aField `env:"omitprefix"`
	}

	type Embedded struct {
		EmbeddedField string
	}

	type EmbeddedPointer struct {
		EmbeddedPointerField string
	}

	type EmbeddedWPrefix struct {
		EmbeddedWPrefixField string
	}

	type EmbeddedCustomPrefix struct {
		EmbeddedCustomField string
	}

	type SliceStruct struct {
		SliceField string
	}

	type CustomIntType int

	type level1 struct {
		// note: private embedded structs are not accessible.
		Embedded             `env:"omitprefix"` // if omitprefix not provided then the prefix is "Embedded"
		*EmbeddedPointer     `env:"omitprefix"` // pointer also is valid
		EmbeddedWPrefix                         // keep prefix
		EmbeddedCustomPrefix `env:"E"`

		DurField                    time.Duration
		TimeField                   time.Time  // default format is RFC3339
		TimeCustomField             time.Time  `fmt:"2006/01/02"` // custom format
		TimePointerField            *time.Time // *time.Time supported but be careful with referencing time.Time!
		FirstField                  *string    `env:"first_field"`
		SecondField                 string
		IntField                    int
		CustomIntField              CustomIntType // just treated as an int
		IntPointerField             *int
		BoolField                   bool
		BoolFieldFalse              bool // default is true but set env is false so final value should be false.
		BoolPointerField            *bool
		ArrayField                  [3]int
		SliceStringField            []string
		SliceIntField               []int
		SliceIntFieldWSpaces        []int // input env value should be able to be '1, 2, 3'
		SliceIntFieldWQuotes1       []int // input env value should be able to be '"1","2","3"'
		SliceIntFieldWQuotes2       []int // input env value should be able to be "'1','2','3''
		SliceIntFieldSquareBrackets []int // input env value should be able to be "[1,2,3]"
		SliceFloatField             []float32
		SliceStructField            []SliceStruct     // slice of structs is ignored
		MapField                    map[string]string // maps are ignored.
		IgnoreField                 string            `env:"-"` // ignore field
		IgnoreStruct                level2            `env:"-"` // ignore struct
		IgnorePointerStruct         *level2           `env:"-"` // ignore struct pointer (will not even be initialized)

		// omitprefix
		// this level omits prefix but the next one does not.
		OmitPrefix        withprefix  `env:"omitprefix"`
		OmitPrefixPointer *withprefix `env:"omitprefix"`

		// omitprefix fallthrough
		// prefix at this level but next level prefix omitted.
		WithPrefixInherited        omitprefix
		WithPrefixInheritedPointer *omitprefix

		Level2 level2 `env:"LEVEL2"`

		privateField     string // should not be populated
		privateFieldWTag string `env:"private_field_w_tag"` // Will not populate private field even with tag.
	}

	cfg := level1{
		BoolFieldFalse: true,
	}

	// error if struct is not a pointer
	d := &Decoder{}
	err := d.Unmarshal(cfg)
	assert.NotNil(t, err)

	// set env vars
	os.Setenv("EMBEDDED_FIELD", "vEMBEDDED_FIELD")
	os.Setenv("EMBEDDED_POINTER_FIELD", "vEMBEDDED_POINTER_FIELD")
	os.Setenv("EMBEDDED_W_PREFIX_EMBEDDED_W_PREFIX_FIELD", "vEMBEDDED_W_PREFIX_EMBEDDED_W_PREFIX_FIELD")
	os.Setenv("E_EMBEDDED_CUSTOM_FIELD", "vE_EMBEDDED_CUSTOM_FIELD")
	os.Setenv("DUR_FIELD", "64s")
	os.Setenv("TIME_FIELD", "2020-01-01T11:04:01Z")         // RFC3339
	os.Setenv("TIME_CUSTOM_FIELD", "2020/01/02")            // custom format
	os.Setenv("TIME_POINTER_FIELD", "2020-01-01T11:04:01Z") // RFC3339
	os.Setenv("first_field", "vfirst_field")
	os.Setenv("SECOND_FIELD", "vSECOND_FIELD")
	os.Setenv("INT_FIELD", "1")
	os.Setenv("CUSTOM_INT_FIELD", "3")
	os.Setenv("INT_POINTER_FIELD", "2")
	os.Setenv("BOOL_FIELD", "true")
	os.Setenv("BOOL_FIELD_FALSE", "false") // false should overwrite default of true
	os.Setenv("BOOL_POINTER_FIELD", "true")
	os.Setenv("ARRAY_FIELD", "1,2,3")
	os.Setenv("SLICE_STRING_FIELD", "part1,part2")
	os.Setenv("SLICE_INT_FIELD", "1,2,3")
	os.Setenv("SLICE_INT_FIELD_W_SPACES", "1, 2, 3")
	os.Setenv("SLICE_INT_FIELD_W_QUOTES_1", `"1","2","3"`)
	os.Setenv("SLICE_INT_FIELD_W_QUOTES_2", `'1','2','3'`)
	os.Setenv("SLICE_INT_FIELD_SQUARE_BRACKETS", "[1,2,3]")
	os.Setenv("SLICE_FLOAT_FIELD", "1.1,2.2,3.3")
	os.Setenv("SLICE_STRUCT_FIELD", "ignored")                                   // structs as single values are ignored.
	os.Setenv("IGNORE_FIELD", "vIGNORE_FIELD")                                   // should not get populated
	os.Setenv("-", "vIGNORE_FIELD")                                              // make sure it doesn't look for a '-' env variable.
	os.Setenv("IGNORE_STRUCT", "vIGNORE_STRUCT")                                 // should not get populated
	os.Setenv("IGNORE_POINTER_STRUCT", "vIGNORE_POINTER_STRUCT")                 // should not get populated
	os.Setenv("WITH_PREFIX_FIELD_1", "vWITH_PREFIX_FIELD_1")                     // field should have this name (top level prefix omitted but next level retained).
	os.Setenv("WITH_PREFIX_INHERITED_FIELD_1", "vWITH_PREFIX_INHERITED_FIELD_1") // top level has prefix but next level ignores it.
	os.Setenv("WITH_PREFIX_INHERITED_POINTER_FIELD_1", "vWITH_PREFIX_INHERITED_POINTER_FIELD_1")
	os.Setenv("PRIVATE_FIELD", "vPRIVATE_FIELD")             // should not get set
	os.Setenv("private_field_w_tag", "vprivate_field_w_tag") // should not get set
	os.Setenv("PRIVATE_FIELD_W_TAG", "vPRIVATE_FIELD_W_TAG") // just checking this variation in case a logic slip.
	os.Setenv("LEVEL2_FIRST_FIELD", "vLEVEL2_FIRST_FIELD")
	os.Setenv("LEVEL2_second_field", "vLEVEL2_second_field")
	os.Setenv("LEVEL2_PRIVATE_FIELD", "vLEVEL2_PRIVATE_FIELD") // should not get set
	os.Setenv("LEVEL2_LEVEL3_FIRST_FIELD", "vLEVEL2_LEVEL3_FIRST_FIELD")
	os.Setenv("LEVEL2_LEVEL3_second_field", "vLEVEL2_LEVEL3_second_field")

	d = &Decoder{}
	cfg = level1{}
	err = d.Unmarshal(&cfg)
	assert.Nil(t, err)

	// make sure each field is populated as expected.
	drtn, _ := time.ParseDuration("64s")
	dte := time.Date(2020, 01, 01, 11, 04, 01, 0, time.UTC)
	assert.Equal(t, cfg.EmbeddedField, "vEMBEDDED_FIELD")
	assert.Equal(t, cfg.EmbeddedPointerField, "vEMBEDDED_POINTER_FIELD")
	assert.Equal(t, cfg.EmbeddedWPrefixField, "vEMBEDDED_W_PREFIX_EMBEDDED_W_PREFIX_FIELD")
	assert.Equal(t, cfg.EmbeddedCustomField, "vE_EMBEDDED_CUSTOM_FIELD")
	assert.Equal(t, cfg.DurField, drtn)
	assert.Equal(t, cfg.TimeField, dte)
	assert.Equal(t, cfg.TimeCustomField, time.Date(2020, 01, 02, 0, 0, 0, 0, time.UTC))
	assert.Equal(t, *cfg.TimePointerField, dte)
	assert.Equal(t, *cfg.FirstField, "vfirst_field")
	assert.Equal(t, cfg.SecondField, "vSECOND_FIELD")
	assert.Equal(t, cfg.IntField, 1)
	assert.Equal(t, int(cfg.CustomIntField), 3) // custom int type just treated as an int.
	assert.Equal(t, *cfg.IntPointerField, 2)
	assert.Equal(t, cfg.BoolField, true)
	assert.Equal(t, cfg.BoolFieldFalse, false)
	assert.Equal(t, *cfg.BoolPointerField, true)
	//assert.Equal(t, cfg.ArrayField, [3]int{1,2,3})
	assert.Equal(t, cfg.SliceStringField, []string{"part1", "part2"})
	assert.Equal(t, cfg.SliceIntField, []int{1, 2, 3})
	assert.Equal(t, cfg.SliceIntFieldWSpaces, []int{1, 2, 3})
	assert.Equal(t, cfg.SliceIntFieldWQuotes1, []int{1, 2, 3})
	assert.Equal(t, cfg.SliceIntFieldSquareBrackets, []int{1, 2, 3})
	assert.Equal(t, cfg.SliceFloatField, []float32{1.1, 2.2, 3.3})
	assert.Empty(t, cfg.IgnoreField)
	assert.Empty(t, cfg.IgnoreStruct)
	assert.Empty(t, cfg.IgnorePointerStruct)
	assert.Equal(t, cfg.OmitPrefix.WithPrefix.Field1, "vWITH_PREFIX_FIELD_1")
	assert.Equal(t, cfg.OmitPrefixPointer.WithPrefix.Field1, "vWITH_PREFIX_FIELD_1")
	assert.Equal(t, cfg.WithPrefixInherited.NoPrefix.Field1, "vWITH_PREFIX_INHERITED_FIELD_1")
	assert.Equal(t, cfg.WithPrefixInheritedPointer.NoPrefix.Field1, "vWITH_PREFIX_INHERITED_POINTER_FIELD_1")
	assert.Empty(t, cfg.privateField)
	assert.Empty(t, cfg.privateFieldWTag)
	assert.Equal(t, *cfg.Level2.FirstField, "vLEVEL2_FIRST_FIELD")
	assert.Equal(t, cfg.Level2.SecondField, "vLEVEL2_second_field")
	assert.Empty(t, cfg.Level2.privateField)
	assert.Equal(t, *cfg.Level2.Level3.FirstField, "vLEVEL2_LEVEL3_FIRST_FIELD")
	assert.Equal(t, cfg.Level2.Level3.SecondField, "vLEVEL2_LEVEL3_second_field")

	// misc tests
	// Test: 'omitprefix' on non-struct and pointer non-struct
	type omitprefixNonStruct struct {
		OmitPrefixField string `env:"omitprefix"` // not allowed returns error.
	}

	d = &Decoder{}
	cfgErr := omitprefixNonStruct{}
	err = d.Unmarshal(&cfgErr)
	assert.EqualError(t, err, "'omitprefix' cannot be used on non-struct field types")

	// Test: a comma in the env tag value gets translated directly as an env field
	// same as everything else. While it doesn't return an error the user is unlikely
	// to set an env variable with a comma. Regardless, the behavior is defined.
	type envComma struct {
		CommaField string `env:"commafield,"`
	}

	os.Setenv("commafield,", "vcommafield,")

	d = &Decoder{}
	cfgComma := envComma{}
	err = d.Unmarshal(&cfgComma)
	assert.Nil(t, err)
	assert.Equal(t, "vcommafield,", cfgComma.CommaField)

	// Test: incorrect formatting - tag value is omitted. only 'env' is provided.
	type envNoValue struct {
		NoTagValueField  string `env:""` // does not return error but has no effect.
		NoTagValueField2 string `env`    // not even the ':""' provided.
	}

	os.Setenv("NO_TAG_VALUE_FIELD", "vNO_TAG_VALUE_FIELD")
	os.Setenv("NO_TAG_VALUE_FIELD_2", "vNO_TAG_VALUE_FIELD_2")

	d = &Decoder{}
	cfgEnvNoValue := envNoValue{}
	err = d.Unmarshal(&cfgEnvNoValue)
	assert.Nil(t, err)
	assert.Equal(t, "vNO_TAG_VALUE_FIELD", cfgEnvNoValue.NoTagValueField)
	assert.Equal(t, "vNO_TAG_VALUE_FIELD_2", cfgEnvNoValue.NoTagValueField2)

	// Test: default values are overwritten.
	// If a default value is provided but no env is found, the default is retained.
	type withDefaults struct {
		DefaultField1 string
		DefaultField2 string
	}

	os.Setenv("DEFAULT_FIELD_1", "vDEFAULT_FIELD_1")

	d = &Decoder{}
	cfgWithDefaults := withDefaults{
		DefaultField1: "default1", // should be overwritten.
		DefaultField2: "default2", // should persist with no env set.
	}
	err = d.Unmarshal(&cfgWithDefaults)
	assert.Nil(t, err)
	assert.Equal(t, "vDEFAULT_FIELD_1", cfgWithDefaults.DefaultField1)
	assert.Equal(t, "default2", cfgWithDefaults.DefaultField2)

	// can only assign "true", "false" or "" to type bool
	type badBool struct {
		BadBoolField bool
	}

	os.Setenv("BAD_BOOL_FIELD", "badvalue") // must be "true", "false", ""

	d = &Decoder{}
	cfgBadBool := badBool{}
	err = d.Unmarshal(&cfgBadBool)
	assert.EqualError(t, err, "'badvalue' from 'BAD_BOOL_FIELD' cannot be set to BadBoolField (bool)")

	// can only assign proper int value to int type.
	type badInt struct {
		BadIntField int
	}

	os.Setenv("BAD_INT_FIELD", "badvalue")

	d = &Decoder{}
	cfgBadInt := badInt{}
	err = d.Unmarshal(&cfgBadInt)
	assert.EqualError(t, err, "'badvalue' from 'BAD_INT_FIELD' cannot be set to BadIntField (int)")

	// test bad uint field
	type badUint struct {
		BadUintField uint
	}

	os.Setenv("BAD_UINT_FIELD", "badvalue")

	d = &Decoder{}
	cfgBadUint := badUint{}
	err = d.Unmarshal(&cfgBadUint)
	assert.EqualError(t, err, "'badvalue' from 'BAD_UINT_FIELD' cannot be set to BadUintField (uint)")

	// Test: pass in pointer of non-struct
	otherPtr := 5

	d = &Decoder{}
	err = d.Unmarshal(&otherPtr)
	assert.EqualError(t, err, "'*int' must be a non-nil pointer struct")

	// teardown: unset envs
	os.Clearenv()
}*/

func TestDecoder_Unmarshal(t *testing.T) {
	type Aint int
	type Astring string
	type tConfig struct {
		Dura   time.Duration
		Time   time.Time `format:"2006-01-02"`
		Bool   bool
		String string

		Int   int
		Int8  int8  `env:"INT8"`
		Int16 int16 `env:"INT16"`
		Int32 int32 `env:"INT32"`
		Int64 int64 `env:"INT64"`

		Uint   uint
		Uint8  uint8  `env:"UINT8"`
		Uint16 uint16 `env:"UINT16"`
		Uint32 uint32 `env:"UINT32"`
		Uint64 uint64 `env:"UINT64"`

		Float32 float32 `env:"FLOAT32"`
		Float64 float64 `env:"FLOAT64"`

		IntP    *int     `env:"INTP"`
		UintP   *uint    `env:"UINTP"`
		StringP *string  `env:"STRINGP"`
		FloatP  *float64 `env:"FLOATP"`

		// slices
		ArrayField       [3]int
		SliceStringField []string
		SliceIntField    []int
		SliceFloatField  []float64

		// ignored fields
		IgnoreMe         string `env:"-"`
		privateField     string // should not be populated
		privateFieldWTag string `env:"private_field_w_tag"` // Will not populate private field even with tag.

		//alias'
		Aint    Aint
		Astring Astring
	}
	type tStruct struct {
		MStruct mStruct  `env:"MSTRUCT"`
		PStruct *mStruct `env:"PSTRUCT"`
	}
	type input struct {
		config interface{}
		args   map[string]string
	}
	fn := func(args ...interface{}) (interface{}, error) {
		os.Clearenv()
		errs := appenderr.New()
		in := args[0].(input)
		if in.config == nil {
			in.config = &tConfig{}
		}
		for key, value := range in.args {
			errs.Add(os.Setenv(key, value))
		}

		d := &Decoder{}
		errs.Add(d.Unmarshal(in.config))
		return in.config, errs.ErrOrNil()
	}
	cases := trial.Cases{
		"default": {
			Input:    input{args: map[string]string{}},
			Expected: &tConfig{},
		},
		"time.duration": {
			Input:    input{args: map[string]string{"DURA": "10s"}},
			Expected: &tConfig{Dura: 10 * time.Second},
		},
		"time.duration (int)": {
			Input:    input{args: map[string]string{"DURA": "1000"}},
			Expected: &tConfig{Dura: 1000},
		},
		"time.Time": {
			Input:    input{args: map[string]string{"TIME": "2010-01-02"}},
			Expected: &tConfig{Time: trial.TimeDay("2010-01-02")},
		},
		"int": {
			Input:    input{args: map[string]string{"INT": "10", "INT8": "8", "INT16": "16", "INT32": "32", "INT64": "64"}},
			Expected: &tConfig{Int: 10, Int8: 8, Int16: 16, Int32: 32, Int64: 64},
		},
		"uint": {
			Input:    input{args: map[string]string{"UINT": "10", "UINT8": "8", "UINT16": "16", "UINT32": "32", "UINT64": "64"}},
			Expected: &tConfig{Uint: 10, Uint8: 8, Uint16: 16, Uint32: 32, Uint64: 64},
		},
		"float": {
			Input:    input{args: map[string]string{"FLOAT32": "3.2", "FLOAT64": "6.4"}},
			Expected: &tConfig{Float32: 3.2, Float64: 6.4},
		},
		"bool=true": {
			Input:    input{args: map[string]string{"BOOL": "true"}},
			Expected: &tConfig{Bool: true},
		},
		"bool=false": {
			Input:    input{args: map[string]string{"BOOL": "false"}},
			Expected: &tConfig{Bool: false},
		},
		"bool (default)": {
			Input:    input{args: map[string]string{}},
			Expected: &tConfig{Bool: false},
		},
		"string": {
			Input:    input{args: map[string]string{"STRING": "hello"}},
			Expected: &tConfig{String: "hello"},
		},
		"array/slice": {
			Input: input{
				args: map[string]string{
					"SLICE_STRING_FIELD": "part1,part2",
					"ARRAY_FIELD":        "1,2,3",
					"SLICE_INT_FIELD":    "1,2,3",
					"SLICE_FLOAT_FIELD":  "1.1,2.2,3.3",
				},
			},
			Expected: &tConfig{
				ArrayField:       [3]int{1, 2, 3},
				SliceStringField: []string{"part1", "part2"},
				SliceIntField:    []int{1, 2, 3},
				SliceFloatField:  []float64{1.1, 2.2, 3.3},
			},
		},
		"slice with quotes\"": {
			Input: input{
				args: map[string]string{"SLICE_INT_FIELD": `"1","2","3"`},
			},
			Expected: &tConfig{
				SliceIntField: []int{1, 2, 3},
			},
		},
		"slice with quotes'": {
			Input: input{
				args: map[string]string{"SLICE_INT_FIELD": `'1','2','3'`},
			},
			Expected: &tConfig{
				SliceIntField: []int{1, 2, 3},
			},
		},
		"slice with brackets[]": {
			Input: input{
				args: map[string]string{"SLICE_INT_FIELD": `[1,2,3]`},
			},
			Expected: &tConfig{
				SliceIntField: []int{1, 2, 3},
			},
		},
		"ignored": {
			Input: input{
				args: map[string]string{
					"IGNORE_ME":           "WHAT?",
					"PRIVATE_FIELD":       "hello",
					"private_field_w_tag": "vprivate_field_w_tag",
					"PRIVATE_FIELD_W_TAG": "vPRIVATE_FIELD_W_TAG", // just checking this variation in case a logic slip.
				},
			},
			Expected: &tConfig{},
		},
		"alias": {
			Input: input{args: map[string]string{"AINT": "10", "ASTRING": "abc", "NUMBER": "two"}},
			Expected: &tConfig{
				Aint:    10,
				Astring: "abc",
			},
		},
		"pointers": {
			Input:    input{args: map[string]string{"INTP": "3", "UINTP": "4", "FLOATP": "5.6"}},
			Expected: &tConfig{IntP: trial.IntP(3), UintP: trial.UintP(4), FloatP: trial.Float64P(5.6)},
		},
		"struct": {
			Input: input{
				config: &tStruct{},
				args:   map[string]string{"MSTRUCT": "abc", "PSTRUCT": "def"},
			},
			Expected: &tStruct{MStruct: mStruct{"abc"}, PStruct: &mStruct{"def"}},
		},
		"keep value for default": {
			Input: input{
				config: &tConfig{
					Int:     1,
					Uint:    2,
					Float64: 3.4,
					String:  "abc",
				},
			},
			Expected: &tConfig{
				Int:     1,
				Uint:    2,
				Float64: 3.4,
				String:  "abc",
			},
		},
	}
	trial.New(fn, cases).SubTest(t)
}

type mStruct struct {
	value string
}

func (m mStruct) MarshalText() ([]byte, error) {
	return []byte(m.value), nil
}

func (m *mStruct) UnmarshalText(b []byte) error {
	m.value = string(b)
	return nil
}
