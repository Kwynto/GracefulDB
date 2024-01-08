package gauth

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gtypes"
	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
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

func Test_blockUser(t *testing.T) {
	Start()

	t.Run("blockUser() function testing - negative", func(t *testing.T) {
		if err := blockUser("fakeuser"); err == nil {
			t.Error("blockUser() error.")
		}
	})

	t.Run("blockUser() function testing - positive", func(t *testing.T) {
		randStr := generateTicket()
		prof := TProfile{
			Description: "Testing description",
			Status:      ACTIVE,
			Roles:       []TRole{USER},
		}
		addUser(randStr, randStr, prof)

		if err := blockUser(randStr); err != nil {
			t.Error("blockUser() error.")
		}
	})

	t.Run("blockUser() function testing - negative", func(t *testing.T) {
		randStr := generateTicket()
		prof := TProfile{
			Description: "Testing description",
			Status:      ACTIVE,
			Roles:       []TRole{USER},
		}
		addUser(randStr, randStr, prof)
		delete(AccessMap, randStr)

		if err := blockUser(randStr); err == nil {
			t.Error("blockUser() error.")
		}
	})
}

func Test_unblockUser(t *testing.T) {
	Start()

	t.Run("unblockUser() function testing - negative", func(t *testing.T) {
		if err := unblockUser("fakeuser"); err == nil {
			t.Error("unblockUser() error.")
		}
	})

	t.Run("unblockUser() function testing - positive", func(t *testing.T) {
		randStr := generateTicket()
		prof := TProfile{
			Description: "Testing description",
			Status:      ACTIVE,
			Roles:       []TRole{USER},
		}
		addUser(randStr, randStr, prof)

		if err := unblockUser(randStr); err != nil {
			t.Error("unblockUser() error.")
		}
	})

	t.Run("unblockUser() function testing - negative", func(t *testing.T) {
		randStr := generateTicket()
		prof := TProfile{
			Description: "Testing description",
			Status:      ACTIVE,
			Roles:       []TRole{USER},
		}
		addUser(randStr, randStr, prof)
		delete(AccessMap, randStr)

		if err := unblockUser(randStr); err == nil {
			t.Error("unblockUser() error.")
		}
	})
}

func Test_updateProfile(t *testing.T) {
	Start()

	t.Run("updateProfile() function testing - negative", func(t *testing.T) {
		prof := TProfile{
			Description: "Testing description",
			Status:      ACTIVE,
			Roles:       []TRole{USER},
		}

		if err := updateProfile("fakeuser", prof); err == nil {
			t.Error("updateProfile() error.")
		}
	})

	t.Run("updateProfile() function testing - positive", func(t *testing.T) {
		randStr := generateTicket()
		prof := TProfile{
			Description: "Testing description",
			Status:      ACTIVE,
			Roles:       []TRole{USER},
		}
		addUser(randStr, randStr, prof)

		if err := updateProfile(randStr, prof); err != nil {
			t.Error("updateProfile() error.")
		}
	})
}

func Test_AddUser(t *testing.T) {
	Start()

	t.Run("AddUser() function testing - negative", func(t *testing.T) {
		prof := TProfile{
			Status: ACTIVE,
			Roles:  []TRole{ADMIN},
		}

		if err := AddUser("root", "toor", prof); err == nil {
			t.Error("AddUser() error.")
		}
	})

	t.Run("AddUser() function testing - positive", func(t *testing.T) {
		randStr := generateTicket()
		prof := TProfile{
			Description: "Testing description",
			Status:      NEW,
			Roles:       []TRole{ADMIN},
		}

		if err := AddUser(randStr, randStr, prof); err != nil {
			t.Error("AddUser() error.")
		}
	})
}

func Test_UpdateUser(t *testing.T) {
	Start()

	t.Run("UpdateUser() function testing - negative", func(t *testing.T) {
		prof := TProfile{
			Status: ACTIVE,
			Roles:  []TRole{ADMIN},
		}

		if err := UpdateUser("fakeuser", "toor", prof); err == nil {
			t.Error("UpdateUser() error.")
		}
	})

	t.Run("UpdateUser() function testing - positive", func(t *testing.T) {
		prof := TProfile{
			Description: "Testing description",
			Status:      ACTIVE,
			Roles:       []TRole{ADMIN},
		}

		if err := UpdateUser("root", "toor", prof); err != nil {
			t.Error("UpdateUser() error.")
		}
	})

	t.Run("UpdateUser() function testing - positive", func(t *testing.T) {
		randStr := generateTicket()
		prof := TProfile{
			Description: "Testing description",
			Status:      ACTIVE,
			Roles:       []TRole{USER},
		}
		addUser(randStr, randStr, prof)

		if err := UpdateUser(randStr, randStr, prof); err != nil {
			t.Error("UpdateUser() error.")
		}
	})
}

func Test_DeleteUser(t *testing.T) {
	Start()

	t.Run("DeleteUser() function testing - negative", func(t *testing.T) {
		if err := DeleteUser("fakeuser"); err == nil {
			t.Error("DeleteUser() error.")
		}
	})

	t.Run("DeleteUser() function testing - positive", func(t *testing.T) {
		randStr := generateTicket()
		prof := TProfile{
			Description: "Testing description",
			Status:      ACTIVE,
			Roles:       []TRole{USER},
		}
		addUser(randStr, randStr, prof)

		if err := DeleteUser(randStr); err != nil {
			t.Error("DeleteUser() error.")
		}
	})

	t.Run("DeleteUser() function testing - positive", func(t *testing.T) {
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

		if err := DeleteUser(randStr); err != nil {
			t.Error("DeleteUser() error.")
		}
	})

	t.Run("DeleteUser() function testing - negative", func(t *testing.T) {
		if err := DeleteUser("root"); err == nil {
			t.Error("DeleteUser() error.")
		}
	})

}

func Test_BlockUser(t *testing.T) {
	Start()

	t.Run("BlockUser() function testing - negative", func(t *testing.T) {
		if err := BlockUser("fakeuser"); err == nil {
			t.Error("BlockUser() error.")
		}
	})

	t.Run("BlockUser() function testing - positive", func(t *testing.T) {
		randStr := generateTicket()
		prof := TProfile{
			Description: "Testing description",
			Status:      ACTIVE,
			Roles:       []TRole{USER},
		}
		addUser(randStr, randStr, prof)

		if err := BlockUser(randStr); err != nil {
			t.Error("BlockUser() error.")
		}
	})

	t.Run("BlockUser() function testing - negative", func(t *testing.T) {
		randStr := generateTicket()
		prof := TProfile{
			Description: "Testing description",
			Status:      ACTIVE,
			Roles:       []TRole{USER},
		}
		addUser(randStr, randStr, prof)
		delete(AccessMap, randStr)

		if err := BlockUser(randStr); err == nil {
			t.Error("BlockUser() error.")
		}
	})

	t.Run("BlockUser() function testing - negative", func(t *testing.T) {
		if err := BlockUser("root"); err == nil {
			t.Error("BlockUser() error.")
		}
	})
}

func Test_UnblockUser(t *testing.T) {
	Start()

	t.Run("UnblockUser() function testing - negative", func(t *testing.T) {
		if err := UnblockUser("fakeuser"); err == nil {
			t.Error("UnblockUser() error.")
		}
	})

	t.Run("UnblockUser() function testing - positive", func(t *testing.T) {
		randStr := generateTicket()
		prof := TProfile{
			Description: "Testing description",
			Status:      ACTIVE,
			Roles:       []TRole{USER},
		}
		addUser(randStr, randStr, prof)

		if err := UnblockUser(randStr); err != nil {
			t.Error("UnblockUser() error.")
		}
	})

	t.Run("UnblockUser() function testing - negative", func(t *testing.T) {
		randStr := generateTicket()
		prof := TProfile{
			Description: "Testing description",
			Status:      ACTIVE,
			Roles:       []TRole{USER},
		}
		addUser(randStr, randStr, prof)
		delete(AccessMap, randStr)

		if err := UnblockUser(randStr); err == nil {
			t.Error("UnblockUser() error.")
		}
	})

	t.Run("UnblockUser() function testing - negative", func(t *testing.T) {
		if err := UnblockUser("root"); err == nil {
			t.Error("UnblockUser() error.")
		}
	})
}

func Test_UpdateProfile(t *testing.T) {
	Start()

	t.Run("UpdateProfile() function testing - negative", func(t *testing.T) {
		prof := TProfile{
			Description: "Testing description",
			Status:      ACTIVE,
			Roles:       []TRole{USER},
		}

		if err := UpdateProfile("fakeuser", prof); err == nil {
			t.Error("UpdateProfile() error.")
		}
	})

	t.Run("UpdateProfile() function testing - positive", func(t *testing.T) {
		randStr := generateTicket()
		prof := TProfile{
			Description: "Testing description",
			Status:      ACTIVE,
			Roles:       []TRole{USER},
		}
		addUser(randStr, randStr, prof)

		if err := UpdateProfile(randStr, prof); err != nil {
			t.Error("UpdateProfile() error.")
		}
	})
}

func Test_CheckUser(t *testing.T) {
	Start()

	t.Run("CheckUser() function testing - negative", func(t *testing.T) {
		if b := CheckUser("fakeuser", "fakeuser"); b {
			t.Error("CheckUser() error.")
		}
	})

	t.Run("CheckUser() function testing - positive", func(t *testing.T) {
		if b := CheckUser("root", "toor"); !b {
			t.Error("CheckUser() error.")
		}
	})

	t.Run("CheckUser() function testing - negative", func(t *testing.T) {
		if b := CheckUser("root", "root"); b {
			t.Error("CheckUser() error.")
		}
	})
}

func Test_GetProfile(t *testing.T) {
	Start()

	t.Run("GetProfile() function testing - negative", func(t *testing.T) {
		if _, err := GetProfile("fakeuser"); err == nil {
			t.Error("GetProfile() error.")
		}
	})

	t.Run("GetProfile() function testing - positive", func(t *testing.T) {
		if _, err := GetProfile("root"); err != nil {
			t.Error("GetProfile() error.")
		}
	})

	t.Run("GetProfile() function testing", func(t *testing.T) {
		res, _ := GetProfile("fakeuser")
		if reflect.TypeOf(res) != reflect.TypeOf(TProfile{}) {
			t.Error("GetProfile() error = The function returns the wrong type")
		}
	})

	t.Run("GetProfile() function testing", func(t *testing.T) {
		res, _ := GetProfile("root")
		if reflect.TypeOf(res) != reflect.TypeOf(TProfile{}) {
			t.Error("GetProfile() error = The function returns the wrong type")
		}
	})
}

func Test_CheckTicket(t *testing.T) {
	Start()

	t.Run("CheckTicket() function testing - positive", func(t *testing.T) {
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
		ticket, _ := NewAuth(&secret)

		_, _, _, err := CheckTicket(ticket)
		if err != nil {
			t.Error("CheckTicket() error.")
		}
	})

	t.Run("CheckTicket() function testing - negative", func(t *testing.T) {
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
		ticket, _ := NewAuth(&secret)
		delete(AccessMap, randStr)

		// login, access, newticket, err := CheckTicket("fakeuser")
		_, _, _, err := CheckTicket(ticket)
		if err == nil {
			t.Error("CheckTicket() error.")
		}
	})

	t.Run("CheckTicket() function testing (old ticket) - positive", func(t *testing.T) {
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
		ticket, _ := NewAuth(&secret)
		NewAuth(&secret)

		_, _, _, err := CheckTicket(ticket)
		if err != nil {
			t.Error("CheckTicket() error.")
		}
	})

	t.Run("CheckTicket() function testing (old ticket) - negative", func(t *testing.T) {
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
		ticket, _ := NewAuth(&secret)
		NewAuth(&secret)
		delete(ticketMap, randStr)

		_, _, _, err := CheckTicket(ticket)
		if err == nil {
			t.Error("CheckTicket() error.")
		}
	})

	t.Run("CheckTicket() function testing (old ticket) - negative", func(t *testing.T) {
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
		ticket, _ := NewAuth(&secret)
		NewAuth(&secret)
		delete(AccessMap, randStr)

		_, _, _, err := CheckTicket(ticket)
		if err == nil {
			t.Error("CheckTicket() error.")
		}
	})

	t.Run("CheckTicket() function testing (old ticket) - negative", func(t *testing.T) {
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
		ticket, _ := NewAuth(&secret)
		NewAuth(&secret)
		delete(reversOldTicketMap, ticket)

		_, _, _, err := CheckTicket(ticket)
		if err == nil {
			t.Error("CheckTicket() error.")
		}
	})
}

func Test_NewAuth(t *testing.T) {
	Start()

	t.Run("NewAuth() function testing", func(t *testing.T) {
		secret := gtypes.VSecret{
			Login:    "",
			Password: "",
		}

		_, err := NewAuth(&secret)
		if err == nil {
			t.Error("NewAuth() error.")
		}
	})

	t.Run("NewAuth() function testing", func(t *testing.T) {
		h := sha256.Sum256([]byte("abracadabra"))
		pass := fmt.Sprintf("%x", h)
		secret := gtypes.VSecret{
			Login:    "abracadabra",
			Password: "abracadabra",
			Hash:     pass,
		}

		_, err := NewAuth(&secret)
		if err == nil {
			t.Error("NewAuth() error.")
		}
	})

	t.Run("NewAuth() function testing", func(t *testing.T) {
		randStr := generateTicket()
		prof := TProfile{
			Description: "Testing description",
			Status:      ACTIVE,
			Roles:       []TRole{USER},
		}
		addUser(randStr, randStr, prof)

		secret := gtypes.VSecret{
			Login:    randStr,
			Password: "abracadabra",
		}

		_, err := NewAuth(&secret)
		if err == nil {
			t.Error("NewAuth() error.")
		}
	})

	t.Run("NewAuth() function testing", func(t *testing.T) {
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

		_, err := NewAuth(&secret)
		if err != nil {
			t.Error("NewAuth() error.")
		}
	})
}

func Test_hashLoad(t *testing.T) {
	t.Run("hashLoad() function testing - positive", func(t *testing.T) {
		hashLoad()
		_, ok := HashMap["root"]

		if !ok {
			t.Error("hashLoad() error.")
		}
	})

	t.Run("hashLoad() function testing - negative", func(t *testing.T) {
		tf := "../../../../config/develop.yaml"
		AuthFile = tf

		hashLoad()
		_, ok := HashMap["root"]
		if !ok {
			t.Error("hashLoad() error.")
		}
	})
}

func Test_hashSave(t *testing.T) {
	t.Run("hashSave() function testing", func(t *testing.T) {
		tf := "./test.json"
		AuthFile = tf
		tempFile, _ := os.OpenFile(AuthFile, os.O_CREATE|os.O_APPEND, os.ModePerm)
		defer os.Remove(AuthFile)
		defer tempFile.Close()

		hashSave()
		if _, err := os.Stat(AuthFile); os.IsNotExist(err) {
			t.Error("hashSave() error.")
		}
	})
}

func Test_accessLoad(t *testing.T) {
	t.Run("accessLoad() function testing - positive", func(t *testing.T) {
		accessLoad()
		_, ok := AccessMap["root"]
		if !ok {
			t.Error("accessLoad() error.")
		}
	})

	t.Run("accessLoad() function testing - negative", func(t *testing.T) {
		tf := "../../../../config/develop.yaml"
		AccessFile = tf

		accessLoad()
		_, ok := AccessMap["root"]
		if !ok {
			t.Error("accessLoad() error.")
		}
	})
}

func Test_accessSave(t *testing.T) {
	t.Run("accessSave() function testing", func(t *testing.T) {
		tf := "./test.json"
		AccessFile = tf
		tempFile, _ := os.OpenFile(AccessFile, os.O_CREATE|os.O_APPEND, os.ModePerm)
		defer os.Remove(AccessFile)
		defer tempFile.Close()

		accessSave()
		if _, err := os.Stat(AccessFile); os.IsNotExist(err) {
			t.Error("accessSave() error.")
		}
	})
}

func Test_Start(t *testing.T) {
	AuthFile = AUTH_FILE
	AccessFile = ACCESS_FILE

	t.Run("Start() function testing", func(t *testing.T) {
		Start()
		_, ok := HashMap["root"]
		if !ok {
			t.Error("Start() error.")
		}
	})
}

func Test_Shutdown(t *testing.T) {
	AuthFile = AUTH_FILE
	AccessFile = ACCESS_FILE
	closer.AddHandler(Shutdown)

	t.Run("Shutdown() function testing", func(t *testing.T) {
		Shutdown(context.Background(), closer.CloseProcs)
		if closer.CloseProcs.Counter != 0 {
			t.Error("Shutdown() error.")
		}
	})
}
