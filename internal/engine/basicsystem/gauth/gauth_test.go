package gauth

import (
	"reflect"
	"testing"

	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gtypes"
)

const (
	TESTING_ITER = 100
)

func Test_TRole_String(t *testing.T) {
	t.Run("String() function testing", func(t *testing.T) {
		res := SYSTEM.String()
		if res != "SYSTEM" {
			t.Error("String() error = wrong result.")
		}
	})

	t.Run("String() function testing", func(t *testing.T) {
		res := ADMIN.String()
		if res != "ADMIN" {
			t.Error("String() error = wrong result.")
		}
	})

	t.Run("String() function testing", func(t *testing.T) {
		res := MANAGER.String()
		if res != "MANAGER" {
			t.Error("String() error = wrong result.")
		}
	})

	t.Run("String() function testing", func(t *testing.T) {
		res := ENGINEER.String()
		if res != "ENGINEER" {
			t.Error("String() error = wrong result.")
		}
	})

	t.Run("String() function testing", func(t *testing.T) {
		res := USER.String()
		if res != "USER" {
			t.Error("String() error = wrong result.")
		}
	})
}

func Test_TRole_IsSystem(t *testing.T) {
	t.Run("IsSystem() function testing", func(t *testing.T) {
		res := SYSTEM.IsSystem()
		if !res {
			t.Error("IsSystem() error = wrong result.")
		}
	})
}

func Test_TRole_IsAdmin(t *testing.T) {
	t.Run("IsAdmin() function testing", func(t *testing.T) {
		res := ADMIN.IsAdmin()
		if !res {
			t.Error("IsAdmin() error = wrong result.")
		}
	})
}

func Test_TRole_IsManager(t *testing.T) {
	t.Run("IsManager() function testing", func(t *testing.T) {
		res := MANAGER.IsManager()
		if !res {
			t.Error("IsManager() error = wrong result.")
		}
	})
}

func Test_TRole_IsEngineer(t *testing.T) {
	t.Run("IsEngineer() function testing", func(t *testing.T) {
		res := ENGINEER.IsEngineer()
		if !res {
			t.Error("IsEngineer() error = wrong result.")
		}
	})
}

func Test_TRole_IsUser(t *testing.T) {
	t.Run("IsUser() function testing", func(t *testing.T) {
		res := USER.IsUser()
		if !res {
			t.Error("IsUser() error = wrong result.")
		}
	})
}

func Test_TRole_IsNotUser(t *testing.T) {
	t.Run("IsNotUser() function testing", func(t *testing.T) {
		res := SYSTEM.IsNotUser()
		if !res {
			t.Error("IsNotUser() error = wrong result.")
		}
	})
}

func Test_TStatus_String(t *testing.T) {
	t.Run("String() function testing", func(t *testing.T) {
		res := UNDEFINED.String()
		if res != "UNDEFINED" {
			t.Error("String() error = wrong result.")
		}
	})

	t.Run("String() function testing", func(t *testing.T) {
		res := NEW.String()
		if res != "NEW" {
			t.Error("String() error = wrong result.")
		}
	})

	t.Run("String() function testing", func(t *testing.T) {
		res := ACTIVE.String()
		if res != "ACTIVE" {
			t.Error("String() error = wrong result.")
		}
	})

	t.Run("String() function testing", func(t *testing.T) {
		res := BANED.String()
		if res != "BANED" {
			t.Error("String() error = wrong result.")
		}
	})
}

func Test_TStatus_IsBad(t *testing.T) {
	t.Run("IsBad() function testing - negative", func(t *testing.T) {
		res := UNDEFINED.IsBad()
		if !res {
			t.Error("IsBad() error = wrong result.")
		}
	})

	t.Run("IsBad() function testing - positive", func(t *testing.T) {
		res := NEW.IsBad()
		if res {
			t.Error("IsBad() error = wrong result.")
		}
	})
}

func Test_TStatus_IsGood(t *testing.T) {
	t.Run("IsGood() function testing - negative", func(t *testing.T) {
		res := NEW.IsGood()
		if !res {
			t.Error("IsGood() error = wrong result.")
		}
	})

	t.Run("IsGood() function testing - positive", func(t *testing.T) {
		res := UNDEFINED.IsGood()
		if res {
			t.Error("IsBad() error = wrong result.")
		}
	})
}

func Test_TProfile_IsAllowed(t *testing.T) {
	t.Run("IsAllowed() function testing - negative", func(t *testing.T) {
		in := TProfile{
			Status: UNDEFINED,
		}
		rule := []TRole{ADMIN}
		res := in.IsAllowed(rule)
		if res {
			t.Error("IsAllowed() error = wrong result.")
		}
	})

	t.Run("IsAllowed() function testing - positive", func(t *testing.T) {
		in := TProfile{
			Status: ACTIVE,
			Roles:  []TRole{ADMIN},
		}
		rule := []TRole{MANAGER}
		res := in.IsAllowed(rule)
		if !res {
			t.Error("IsAllowed() error = wrong result.")
		}
	})

	t.Run("IsAllowed() function testing - positive", func(t *testing.T) {
		in := TProfile{
			Status: ACTIVE,
			Roles:  []TRole{MANAGER},
		}
		rule := []TRole{ENGINEER, MANAGER}
		res := in.IsAllowed(rule)
		if !res {
			t.Error("IsAllowed() error = wrong result.")
		}
	})

	t.Run("IsAllowed() function testing - negative", func(t *testing.T) {
		in := TProfile{
			Status: ACTIVE,
			Roles:  []TRole{ENGINEER},
		}
		rule := []TRole{MANAGER}
		res := in.IsAllowed(rule)
		if res {
			t.Error("IsAllowed() error = wrong result.")
		}
	})
}

func Test_generateTicket(t *testing.T) {
	t.Run("generateTicket() function testing", func(t *testing.T) {
		etalon := "string"
		res := generateTicket()
		if reflect.TypeOf(res) != reflect.TypeOf(etalon) {
			t.Error("generateTicket() error = The function returns the wrong type")
		}
	})

	t.Run("generateTicket() function testing", func(t *testing.T) {
		testVar := make(map[int]string)
		for i := 0; i < TESTING_ITER; i++ {
			testVar[i] = generateTicket() // calling the tested function
		}
		for _, v1 := range testVar {
			count := 0
			for _, v2 := range testVar {
				if v1 == v2 {
					count++
				}
			}
			// work check
			if count > 1 {
				t.Error("Error generating unique ticket.")
			}
		}
	})
}

func Test_addUser(t *testing.T) {
	Start()

	t.Run("addUser() function testing - negative", func(t *testing.T) {
		prof := TProfile{
			Status: ACTIVE,
			Roles:  []TRole{ADMIN},
		}

		if err := addUser("root", "toor", prof); err == nil {
			t.Error("addUser() error.")
		}
	})

	t.Run("addUser() function testing - positive", func(t *testing.T) {
		randStr := generateTicket()
		prof := TProfile{
			Description: "Testing description",
			Status:      NEW,
			Roles:       []TRole{ADMIN},
		}

		if err := addUser(randStr, randStr, prof); err != nil {
			t.Error("addUser() error.")
		}
	})
}

func Test_updateUser(t *testing.T) {
	Start()

	t.Run("updateUser() function testing - negative", func(t *testing.T) {
		prof := TProfile{
			Status: ACTIVE,
			Roles:  []TRole{ADMIN},
		}

		if err := updateUser("fakeuser", "toor", prof); err == nil {
			t.Error("updateUser() error.")
		}
	})

	t.Run("updateUser() function testing - positive", func(t *testing.T) {
		prof := TProfile{
			Description: "Testing description",
			Status:      ACTIVE,
			Roles:       []TRole{ADMIN},
		}

		if err := updateUser("root", "toor", prof); err != nil {
			t.Error("updateUser() error.")
		}
	})

	t.Run("updateUser() function testing - positive", func(t *testing.T) {
		randStr := generateTicket()
		prof := TProfile{
			Description: "Testing description",
			Status:      ACTIVE,
			Roles:       []TRole{USER},
		}
		addUser(randStr, randStr, prof)

		if err := updateUser(randStr, randStr, prof); err != nil {
			t.Error("updateUser() error.")
		}
	})
}

func Test_deleteUser(t *testing.T) {
	Start()

	t.Run("deleteUser() function testing - negative", func(t *testing.T) {
		if err := deleteUser("fakeuser"); err == nil {
			t.Error("deleteUser() error.")
		}
	})

	t.Run("deleteUser() function testing - positive", func(t *testing.T) {
		randStr := generateTicket()
		prof := TProfile{
			Description: "Testing description",
			Status:      ACTIVE,
			Roles:       []TRole{USER},
		}
		addUser(randStr, randStr, prof)

		if err := deleteUser(randStr); err != nil {
			t.Error("deleteUser() error.")
		}
	})

	t.Run("deleteUser() function testing - positive", func(t *testing.T) {
		randStr := generateTicket()
		prof := TProfile{
			Description: "Testing description",
			Status:      ACTIVE,
			Roles:       []TRole{USER},
		}
		addUser(randStr, randStr, prof)

		secret := gtypes.VSecret{
			Login:    randStr,
			Password: randStr,
		}
		NewAuth(&secret)
		NewAuth(&secret)

		if err := deleteUser(randStr); err != nil {
			t.Error("deleteUser() error.")
		}
	})
}
