package test

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"log"
	"slices"
	"strconv"
	"testing"
)

func TestValidationStruct(t *testing.T) {
	validate := validator.New()

	type Customer struct {
		Nama string `json:"nama,omitempty" validate:"required,min=2"`
	}

	scenario := []struct {
		Name             string
		Input            Customer
		ExpectedErrorNil bool
	}{
		{
			Name:             "test success validation",
			Input:            Customer{"reo"},
			ExpectedErrorNil: true,
		},
		{
			Name:             "test failed validation",
			Input:            Customer{Nama: "n"},
			ExpectedErrorNil: false,
		},
	}

	for _, testScenario := range scenario {
		ctx := context.Background()
		t.Run(testScenario.Name, func(t *testing.T) {
			err := validate.StructCtx(ctx, testScenario.Input)

			assert.Equal(t, err == nil, testScenario.ExpectedErrorNil)
		})
	}
}

func TestValidateVariable(t *testing.T) {
	validate := validator.New()

	scenario := []struct {
		Name          string
		Input         string
		Tag           string
		ExpectedError bool
	}{
		{
			Name:          "test required failed",
			Input:         "",
			Tag:           "required",
			ExpectedError: true,
		},
		{
			Name:          "test required success validate",
			Input:         "reo",
			Tag:           "required",
			ExpectedError: false,
		},
	}

	for _, scenarioTest := range scenario {
		t.Run(scenarioTest.Name, func(t *testing.T) {
			ctx := context.Background()
			err := validate.VarCtx(ctx, scenarioTest.Input, scenarioTest.Tag)
			if err != nil {
				errorField := err.(validator.ValidationErrors)
				log.Printf("error with tag [%v]", errorField[0].Tag())
			}

			assert.Equal(t, err != nil, scenarioTest.ExpectedError)
		})
	}
}

// TestValidasiDuaVariable untuk validasi pada 2 variabel
// misal saat kita melakukan validasi password dan confirmPassword
// contoh : validate.VarWithValue(password, confirmPassword, "eqfield")
func TestValidasiDuaVariable(t *testing.T) {
	validate := validator.New()

	scenario := []struct {
		Name            string
		Password        string
		ConfirmPassword string
		ExpectError     bool
	}{
		{
			Name:            "test validation password failed",
			Password:        "123456",
			ConfirmPassword: "123",
			ExpectError:     true,
		},
		{
			Name:            "test validation password success",
			Password:        "123456",
			ConfirmPassword: "123456",
			ExpectError:     false,
		},
	}

	for _, scTest := range scenario {
		t.Run(scTest.Name, func(t *testing.T) {
			err := validate.VarWithValueCtx(context.Background(), scTest.Password, scTest.ConfirmPassword, "eqfield")
			if err != nil {
				if validationErrors, ok := err.(validator.ValidationErrors); ok {
					for _, errorField := range validationErrors {
						log.Printf("error on field [%v] with tag [%v]", errorField.Field(), errorField.Tag())
					}
				}
			}

			assert.Equal(t, err != nil, scTest.ExpectError)
		})
	}
}

// TestBackedInValidation untuk validasi menggunakan tag bawaan
// test menggunakan tag yang sudah disediakan oleh package validator
// kita hanya perlu menggunakan nama tag nya saja di validate
func TestBackedInValidation(t *testing.T) {
	validate := validator.New()

	scenario := []struct {
		Name        string `json:"name,omitempty"`
		Input       string `json:"input,omitempty"`
		Tag         string `json:"tag,omitempty"`
		ExpectError bool   `json:"expect_error,omitempty"`
	}{
		{
			Name:        "test validation required failed",
			Input:       "",
			Tag:         "required",
			ExpectError: true,
		},
		{
			Name:        "test validation required success",
			Input:       "hello world",
			Tag:         "required",
			ExpectError: false,
		},
		{
			Name:        "test validation patter ip address failed",
			Input:       "172.www",
			Tag:         "ip",
			ExpectError: true,
		},
		{
			Name:        "test validation pattern ip address success",
			Input:       "172.18.41.238",
			Tag:         "ip",
			ExpectError: false,
		},
	}

	for _, scenarioTest := range scenario {
		t.Run(scenarioTest.Name, func(t *testing.T) {
			err := validate.VarCtx(context.Background(), scenarioTest.Input, scenarioTest.Tag)
			if err != nil {
				if validationErrors, ok := err.(validator.ValidationErrors); ok {
					for _, errorField := range validationErrors {
						log.Printf("error with tag [%v]", errorField.Tag())
					}
				}
			}

			assert.Equal(t, err != nil, scenarioTest.ExpectError)
		})
	}
}

// TestMultipleTagValidation untuk validasi menggunakan lebih dari satu tag
// contoh : validate.VarCtx(ctx, input, "required,min=2,max=10")
func TestMultipleTagValidation(t *testing.T) {
	validate := validator.New()

	scenario := []struct {
		Name        string
		Input       string
		Tag         string
		ExpectError bool
	}{
		{
			Name:        "test required and min failed",
			Input:       "reo",
			Tag:         "required,min=3,ip",
			ExpectError: true,
		},
		{
			Name:        "test required and min success",
			Input:       "172.18.231.248",
			Tag:         "required,min=3,ip",
			ExpectError: false,
		},
	}

	for _, scenarioTest := range scenario {
		t.Run(scenarioTest.Name, func(t *testing.T) {
			err := validate.VarCtx(context.Background(), scenarioTest.Input, scenarioTest.Tag)
			if err != nil {
				if validationErrors, ok := err.(validator.ValidationErrors); ok {
					for _, errorField := range validationErrors {
						log.Printf("error with tag [%v] with error [%v]", errorField.Tag(), errorField.Error())
					}
				}
			}

			assert.Equal(t, err != nil, scenarioTest.ExpectError)
		})
	}
}

// TestTagParameter untuk validasi yang perlu menggunakan parameter dalam tag
// misal ingin validasi teks harus minimal 5 karakter, maka menggunakan min=5
// contoh "required,min=5"
func TestTagParameter(t *testing.T) {
	validate := validator.New()

	scenario := []struct {
		Name        string
		Input       string
		Tag         string
		ExpectError bool
	}{
		{
			Name:        "test tag parameter min and max failed",
			Input:       "re",
			Tag:         "required,min=3,max=10",
			ExpectError: true,
		},
		{
			Name:        "test tag parameter min and max success",
			Input:       "reo s",
			Tag:         "required,min=3,max=10",
			ExpectError: false,
		},
	}

	for _, scenarioTest := range scenario {
		t.Run(scenarioTest.Name, func(t *testing.T) {
			err := validate.VarCtx(context.Background(), scenarioTest.Input, scenarioTest.Tag)
			if err != nil {
				if validationErrors, ok := err.(validator.ValidationErrors); ok {
					for _, errorField := range validationErrors {
						log.Printf("error with tag [%v]", errorField.Tag())
					}
				}
			}

			assert.Equal(t, err != nil, scenarioTest.ExpectError)
		})
	}
}

// TestValidasiStruct untuk melakukan validasi pada variabel yang tipe datanya struct
// tuliskan tag pada masing-masing field yang ada di dalam struct
// tambahkan tag `validate:"required"` pada fieldnya
func TestValidasiStruct(t *testing.T) {
	validate := validator.New()

	type LoginRequest struct {
		Username string `json:"username,omitempty" validate:"required,email"`
		Password string `json:"password,omitempty" validate:"required,min=6"`
	}

	scenario := []struct {
		Name          string
		Input         LoginRequest
		ExpectedError bool
	}{
		{
			Name: "test validasi struct failed",
			Input: LoginRequest{
				Username: "reo",
				Password: "123",
			},
			ExpectedError: true,
		},
		{
			Name: "test validasi struct success",
			Input: LoginRequest{
				Username: "reo123@gmail.com",
				Password: "123456",
			},
			ExpectedError: false,
		},
	}

	for _, scenarioTest := range scenario {
		t.Run(scenarioTest.Name, func(t *testing.T) {
			err := validate.StructCtx(context.Background(), scenarioTest.Input)

			assert.Equal(t, err != nil, scenarioTest.ExpectedError)
		})
	}
}

// TestValidationErrors untuk mengecek error yang dikembalikan
func TestValidationErrors(t *testing.T) {
	validate := validator.New()

	type LoginRequet struct {
		Username string `json:"username,omitempty" validate:"required,email,min=3"`
		Password string `json:"password,omitempty" validate:"required,alpha,min=6"`
	}

	scenario := []struct {
		Name        string       `json:"name,omitempty"`
		Input       *LoginRequet `json:"input,omitempty"`
		ExpectError bool         `json:"expect_error,omitempty"`
	}{
		{
			Name: "test struct validation errors failed",
			Input: &LoginRequet{
				Username: "r",
				Password: "123",
			},
			ExpectError: true,
		},
		{
			Name: "test struct validation errors success",
			Input: &LoginRequet{
				Username: "emailtest@gmail.com",
				Password: "qwerty",
			},
			ExpectError: false,
		},
	}

	for _, scenarioTest := range scenario {
		t.Run(scenarioTest.Name, func(t *testing.T) {
			err := validate.StructCtx(context.Background(), *scenarioTest.Input)
			if err != nil {
				if validationErrors, ok := err.(validator.ValidationErrors); ok {
					for _, errorField := range validationErrors {
						log.Printf("error on field [%v] with tag [%v]\n", errorField.Field(), errorField.Tag())
					}
				}
			}

			assert.Equal(t, err != nil, scenarioTest.ExpectError)
		})
	}
}

// TestValidasiNestedStruct untuk melakukan validasi pasa struct dengan field tipe data struct
// misal Address *Address `validate:"required"`
func TestValidasiNestedStruct(t *testing.T) {
	validate := validator.New()

	type Address struct {
		City    string `json:"city,omitempty" validate:"required"`
		Country string `json:"country,omitempty" validate:"required"`
	}

	type User struct {
		Name    string   `json:"name,omitempty" validate:"required"`
		Address *Address `json:"address,omitempty" validate:"required"`
	}

	scenario := []struct {
		Name        string
		Input       *User
		ExpectError bool
	}{
		{
			Name: "test validate nested failed",
			Input: &User{
				Name:    "reo",
				Address: &Address{},
			},
			ExpectError: true,
		},
		{
			Name: "test validate nested success",
			Input: &User{
				Name: "reo sahobby",
				Address: &Address{
					City:    "Jakarta Selatan",
					Country: "Indonesia",
				},
			},
			ExpectError: false,
		},
	}

	for _, scenarioTest := range scenario {
		t.Run(scenarioTest.Name, func(t *testing.T) {
			err := validate.StructCtx(context.Background(), scenarioTest.Input)
			if err != nil {
				if validationErrrors, ok := err.(validator.ValidationErrors); ok {
					for _, errorField := range validationErrrors {
						log.Printf("%v\n", errorField.Error())
					}
				}
			}

			assert.Equal(t, err != nil, scenarioTest.ExpectError)
		})
	}
}

// TestValidasiSlice untuk melakukan validasi tipe data slice/array
// menggunakan tag 'dive'
// contoh Address []Address `validate:"required,dive"`
func TestValidasiSlice(t *testing.T) {
	validate := validator.New()

	type Address struct {
		City    string `json:"city,omitempty" validate:"required,min=2"`
		Country string `json:"country,omitempty" validate:"required,min=2"`
	}

	type User struct {
		Name      string    `json:"name,omitempty" validate:"required"`
		Addresses []Address `json:"addresses,omitempty" validate:"required,dive"`
	}

	scenario := []struct {
		Name        string
		Input       *User
		ExpectError bool
	}{
		{
			Name: "test validate slice failed",
			Input: &User{
				Name: "reo",
				Addresses: []Address{
					{
						City:    "",
						Country: "",
					},
					{
						City:    "",
						Country: "",
					},
				},
			},
			ExpectError: true,
		},
		{
			Name: "test validate slice struct success",
			Input: &User{
				Name: "reo sahobby",
				Addresses: []Address{
					{
						City:    "Jakarta Selatan",
						Country: "Indonesia",
					},
				},
			},
			ExpectError: false,
		},
	}

	for _, testScenario := range scenario {
		t.Run(testScenario.Name, func(t *testing.T) {
			err := validate.StructCtx(context.Background(), testScenario.Input)
			if err != nil {
				for _, errorField := range err.(validator.ValidationErrors) {
					log.Println(errorField.Error())
				}
			}

			assert.Equal(t, err != nil, testScenario.ExpectError)
		})
	}
}

// TestValidasiBasicSlice untuk validasi slice tipe data basic
// contohnya validasi []string `validate:"required,dive,ip"`
// tag validasi valuenya ditulis setelah dive
func TestValidasiBasicSlice(t *testing.T) {
	validate := validator.New()

	type Server struct {
		Name        string   `json:"name,omitempty" validate:"required"`
		IPAddresses []string `json:"ip_addresses,omitempty" validate:"required,dive,ip"`
	}

	scenario := []struct {
		Name        string
		Input       *Server
		ExpectError bool
	}{
		{
			Name: "test validasi basic slice failed",
			Input: &Server{
				Name:        "server dev",
				IPAddresses: []string{"qwert", "wasd"},
			},
			ExpectError: true,
		},
		{
			Name: "test validasi basic slice success",
			Input: &Server{
				Name:        "server production",
				IPAddresses: []string{"172.100.100.101", "172.100.100.102"},
			},
			ExpectError: false,
		},
	}

	for _, testScenario := range scenario {
		t.Run(testScenario.Name, func(t *testing.T) {
			err := validate.StructCtx(context.Background(), testScenario.Input)
			if err != nil {
				for _, errorField := range err.(validator.ValidationErrors) {
					log.Println(errorField.Error())
				}
			}

			assert.Equal(t, err != nil, testScenario.ExpectError)
		})
	}

}

// TestValidasiMap untuk validasi pada tipe data map
// dalam map karea ada key-value. maka kita bisa menambahkan dive untuk key dan value
// menggunakan keys dan endkeys
func TestValidasiMap(t *testing.T) {
	validate := validator.New()

	type School struct {
		Name    string `json:"name,omitempty" validate:"required,min=2"`
		Address string `json:"address,omitempty" validate:"required,min=2"`
	}

	scenario := []struct {
		Name        string
		Input       map[string]*School // variabel yang akan divalidasi
		Tag         string             // tag validation untuk map
		ExpectError bool
	}{
		{
			Name: "test validation map failed",
			Input: map[string]*School{
				"s": &School{
					Name:    "a",
					Address: "a",
				},
			},
			Tag:         "dive,keys,min=2,endkeys,required",
			ExpectError: true,
		},
		{
			Name: "test validatio map success",
			Input: map[string]*School{
				"sd": {
					Name:    "SD N 1",
					Address: "Jakarta Selatan",
				},
			},
			Tag:         "required,dive,keys,min=2,endkeys,required",
			ExpectError: false,
		},
	}

	for _, testScenario := range scenario {
		t.Run(testScenario.Name, func(t *testing.T) {
			err := validate.VarCtx(context.Background(), testScenario.Input, testScenario.Tag)
			if err != nil {
				for _, errorField := range err.(validator.ValidationErrors) {
					log.Println(errorField.Error())
				}
			}

			assert.Equal(t, err != nil, testScenario.ExpectError)
		})
	}
}

// TestValidasiBasicMap untuk validasi pada tipe data basic map[string]string
// karena map punya key-value, untuk vlaidasi tag pada key wajib ditambahkan keys dan endkeys
// contoh : map[string]string `required,dive,keys,required,alpha,endkeys,dive,email,min=5`
func TestValidasiBasicMap(t *testing.T) {
	validate := validator.New()

	scenario := []struct {
		Name        string
		Input       map[string]string // variabel yang akan divalidasi
		Tag         string            // tag validasi untuk variabel map
		ExpectError bool
	}{
		{
			Name: "test validasi basic map failed",
			Input: map[string]string{
				"user1": "user1",
				"user2": "user2@gmail.com",
				"user3": "",
				"a":     "reo@gmail.com",
			},
			Tag:         "required,dive,keys,required,min=3,endkeys,required,email,min=12",
			ExpectError: true,
		},
		{
			Name: "test validasi basic map success",
			Input: map[string]string{
				"server1": "172.18.10.22",
				"server2": "172.18.10.23",
				"server3": "172.18.10.24",
			},
			Tag:         "required,dive,keys,required,endkeys,required,ip",
			ExpectError: false,
		},
	}

	for _, testScenario := range scenario {
		t.Run(testScenario.Name, func(t *testing.T) {
			err := validate.VarCtx(context.Background(), testScenario.Input, testScenario.Tag)
			if err != nil {
				for _, errorField := range err.(validator.ValidationErrors) {
					log.Println(errorField.Error())
				}
			}

			assert.Equal(t, err != nil, testScenario.ExpectError)
		})
	}
}

// TestAliasTag untuk mengganti nama tag sesuai dengan custom tag kita
// valdate.RegisterAlias(alias, tag)
func TestAliasTag(t *testing.T) {
	validate := validator.New()
	validate.RegisterAlias("app_email", "required,email,min=15")

	scenario := []struct {
		Name        string
		Input       string
		ExpectError bool
	}{
		{
			Name:        "test validasi using alias failed",
			Input:       "reoo",
			ExpectError: true,
		},
		{
			Name:        "test validasi using alias success",
			Input:       "reoshby1299@gmail.com",
			ExpectError: false,
		},
	}

	for _, testScenario := range scenario {
		t.Run(testScenario.Name, func(t *testing.T) {
			err := validate.VarCtx(context.Background(), testScenario.Input, "app_email")
			if err != nil {
				for _, errorField := range err.(validator.ValidationErrors) {
					log.Println(errorField.Error())
				}
			}

			assert.Equal(t, err != nil, testScenario.ExpectError)
		})
	}
}

// TestCustomValidation untuk menambahkan costum logic validasi
// misal kita memiliki validasi yang lebih kompleks, bisa buat sendiri logic validasinya
// kemudian kita register aliasnya dengan tag yang kita inginkan
func TestCustomValidation(t *testing.T) {
	validate := validator.New()

	// create constum function validation
	validateType := func(field validator.FieldLevel) bool {
		ourType := []string{"hobby", "gadget", "adventure", "automotive"}

		value, ok := field.Field().Interface().(string)
		if !ok {
			return false
		}
		return slices.Contains(ourType, value)
	}

	// register
	validate.RegisterValidation("category", validateType, false)

	scenario := []struct {
		Name        string
		Input       string
		ExpectError bool
	}{
		{
			Name:        "test validation category failed",
			Input:       "abcd",
			ExpectError: true,
		},
		{
			Name:        "test validation category success",
			Input:       "gadget",
			ExpectError: false,
		},
	}

	for _, testScenario := range scenario {
		t.Run(testScenario.Name, func(t *testing.T) {
			err := validate.VarCtx(context.Background(), testScenario.Input, "category")
			if err != nil {
				for _, errorField := range err.(validator.ValidationErrors) {
					log.Println(errorField.Error())
				}
			}

			assert.Equal(t, err != nil, testScenario.ExpectError)
		})
	}
}

// TestCustomValidationParameter untuk menambahkan costum validasi yang memerlukan parameter
func TestCustomValidationParameter(t *testing.T) {
	validate := validator.New()

	// create function length array
	validateMinCategory := func(field validator.FieldLevel) bool {
		length, err := strconv.Atoi(field.Param())
		if err != nil {
			panic(err)
			return false
		}

		ourType := []string{"a", "b", "c", "d", "e"}
		value, ok := field.Field().Interface().([]string)
		if !ok {
			return false
		}

		// cek eaech type
		for _, catType := range value {
			if !slices.Contains(ourType, catType) {
				return false
			}
		}

		if len(value) < length {
			return false
		}

		return true
	}

	// register validate
	validate.RegisterValidation("min_category", validateMinCategory)

	// create scenario
	scenario := []struct {
		Name        string
		Input       []string
		ExpectError bool
	}{
		{
			Name:        "test validasi category param failed",
			Input:       []string{"a"},
			ExpectError: true,
		},
		{
			Name:        "test validasi category param success",
			Input:       []string{"a", "b", "c"},
			ExpectError: false,
		},
		{
			Name:        "test validasi category not in options",
			Input:       []string{"r", "e", "o"},
			ExpectError: true,
		},
	}

	for _, testScenario := range scenario {
		t.Run(testScenario.Name, func(t *testing.T) {
			err := validate.VarCtx(context.Background(), testScenario.Input, "min_category=2")
			if err != nil {
				for _, errorField := range err.(validator.ValidationErrors) {
					log.Println(errorField.Error())
				}
			}

			assert.Equal(t, err != nil, testScenario.ExpectError)
		})
	}
}

// TestCustomMessageValidation untuk memberi message validasi sesuai dengan custom kita
func TestCustomMessageValidation(t *testing.T) {
	validate := validator.New()

	// register custom gender validation
	genderValidation := func(field validator.FieldLevel) bool {
		value, ok := field.Field().Interface().(string)
		if !ok {
			return false
		}

		return slices.Contains([]string{"male", "female"}, value)
	}

	// add register
	validate.RegisterValidation("gender", genderValidation)

	// add custom Message
	ErrorMessage := func(tag string) string {
		switch tag {
		case "gender":
			return "gender must be male or female"
		default:
			return "unknown error"
		}
	}

	// create scenario
	scenario := []struct {
		Name        string
		Input       string
		ExpectError bool
	}{
		{
			Name:        "test validation gender failed",
			Input:       "mafale",
			ExpectError: true,
		},
		{
			Name:        "test validation gender success",
			Input:       "male",
			ExpectError: false,
		},
	}

	for _, testScenario := range scenario {
		t.Run(testScenario.Name, func(t *testing.T) {
			err := validate.VarCtx(context.Background(), testScenario.Input, "gender")
			if err != nil {
				for _, errorField := range err.(validator.ValidationErrors) {
					log.Println(ErrorMessage(errorField.Tag()))
				}
			}

			assert.Equal(t, err != nil, testScenario.ExpectError)
		})
	}

}
